package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"golang.org/x/time/rate"

	"github.com/mk6i/retro-aim-server/wire"

	"github.com/mk6i/smarter-smarter-child/config"
)

// chatContext stores context for a conversation with a single user.
type chatContext struct {
	// cookie is the unique chat identifier generated by the AIM client.
	cookie uint64
	// lastExchange holds the last two exchanged messages in a chat session to
	// provide contextual memory for ChatGPT, facilitating more relevant and
	// coherent responses. The array format is:
	// [0] - the most recent message received from the user,
	// [1] - the most recent message sent by the bot.
	lastExchange *[2]string
	// limiter enforces rate limits on messages to prevent spam.
	limiter *rate.Limiter
	// rateLimited flags whether the current chat session is being rate limited
	// due to excessive message frequency.
	rateLimited bool
	// semaphore provides synchronization to ensure only 1 response for user is
	// in-flight at any given time.
	semaphore chan struct{}
	// warnCount indicates how many times the user has warned the bot.
	warnCount int
}

func (c chatContext) tryLock() bool {
	select {
	case c.semaphore <- struct{}{}:
		return true
	default:
		return false
	}
}

func (c chatContext) releaseLock() {
	<-c.semaphore
}

// Chat handles conversations with multiple users.
func Chat(logger *slog.Logger, flapc FlapClient, authCookie string, chatBot ChatBot, config config.Config) error {
	if _, err := flapc.ReceiveSignonFrame(); err != nil {
		return err
	}

	tlv := []wire.TLV{
		wire.NewTLV(wire.OServiceTLVTagsLoginCookie, []byte(authCookie)),
	}
	if err := flapc.SendSignonFrame(tlv); err != nil {
		return err
	}

	hostOnlineFrame := wire.SNACFrame{}
	hostOnlineSNAC := wire.SNAC_0x01_0x03_OServiceHostOnline{}
	if err := flapc.ReceiveSNAC(&hostOnlineFrame, &hostOnlineSNAC); err != nil {
		return err
	}

	clientOnlineFrame := wire.SNACFrame{
		FoodGroup: wire.OService,
		SubGroup:  wire.OServiceClientOnline,
	}
	clientOnlineSNAC := wire.SNAC_0x01_0x02_OServiceClientOnline{}
	if err := flapc.SendSNAC(clientOnlineFrame, clientOnlineSNAC); err != nil {
		return err
	}

	msgCh := make(chan wire.SNACMessage, 10)

	// send client->server messages
	go sendSNACs(logger, flapc, msgCh)

	// send heartbeats to the server to keep the connection alive
	go sendHeartbeat(msgCh)

	// keep track of all chat contexts per screen name
	chatContexts := make(map[string]*chatContext)

	logger.Info("listening for incoming IMs")

	for {
		flap, flapBody, err := flapc.ReceiveFLAP()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		if flap.FrameType == wire.FLAPFrameSignoff {
			break // server politely asked us to disconnect
		}
		if flap.FrameType != wire.FLAPFrameData {
			continue // received a non-data FLAP frame, nothing to do here
		}

		snacFrame := wire.SNACFrame{}
		if err := wire.Unmarshal(&snacFrame, flapBody); err != nil {
			return err
		}

		logger.Debug("received SNAC", slog.Group(
			"snac",
			slog.String("foodgroup", wire.FoodGroupName(snacFrame.FoodGroup)),
			slog.String("subgroup", wire.SubGroupName(snacFrame.FoodGroup, snacFrame.SubGroup)),
		))

		switch {
		case snacFrame.FoodGroup == wire.ICBM && snacFrame.SubGroup == wire.ICBMChannelMsgToClient:
			// received an IM, let's respond
			if err := exchangeMessages(logger, msgCh, flapBody, chatContexts, chatBot, config); err != nil {
				return err
			}
		case snacFrame.FoodGroup == wire.OService && snacFrame.SubGroup == wire.OServiceEvilNotification:
			// received a warning, let's respond
			if err := reactToWarning(logger, msgCh, chatContexts, flapBody, chatBot); err != nil {
				return err
			}
		}
	}

	return nil
}

func sendHeartbeat(msgCh chan wire.SNACMessage) {
	for {
		select {
		case <-time.Tick(1 * time.Minute):
			msgCh <- wire.SNACMessage{
				Frame: wire.SNACFrame{
					FoodGroup: wire.OService,
					SubGroup:  wire.OServiceNoop,
				},
				Body: struct{}{},
			}
		}
	}
}

func sendSNACs(logger *slog.Logger, flapc FlapClient, msgCh chan wire.SNACMessage) {
	for msgSNAC := range msgCh {
		group := slog.Group(
			"snac",
			slog.String("foodgroup", wire.FoodGroupName(msgSNAC.Frame.FoodGroup)),
			slog.String("subgroup", wire.SubGroupName(msgSNAC.Frame.FoodGroup, msgSNAC.Frame.SubGroup)),
		)
		if err := flapc.SendSNAC(msgSNAC.Frame, msgSNAC.Body); err != nil {
			logger.Error("error sending SNAC", group)
			return
		}
		logger.Debug("sent SNAC", group)
	}
}

func reactToWarning(
	logger *slog.Logger,
	msgCh chan<- wire.SNACMessage,
	chatContexts map[string]*chatContext,
	flapBody *bytes.Buffer,
	chatBot ChatBot,
) error {

	chatMsg := wire.SNAC_0x01_0x10_OServiceEvilNotification{}
	// io.EOF can be due to an empty OServiceEvilNotification SNAC
	// which indicates an anonymous warning, so ignore it
	if err := wire.Unmarshal(&chatMsg, flapBody); err != nil && err != io.EOF {
		return err
	}

	chatCtx, ok := chatContexts[chatMsg.ScreenName]
	// chatMsg.ScreenName is "" (anonymous), or hasn't sent us an IM yet
	if !ok {
		logger.Debug("can't find chat context, moving on")
		return nil
	}

	chatCtx.warnCount++

	var userMessage string
	switch chatCtx.warnCount {
	case 1:
		userMessage = "Respond in a pleasant tone to me warning you for the first time."
	case 2:
		userMessage = "Respond in a miffed tone to me warning you a second time."
	case 3:
		userMessage = "Respond in an angry tone to me warning you a third time. You are going warn me back right now in retaliation."
	default:
		userMessage = "Respond in an outraged tone to me warning you a fourth time."
	}

	botResponse, err := chatBot.ExchangeMessage(userMessage, *chatCtx.lastExchange)
	if err != nil {
		return fmt.Errorf("unable to get response from bot: %w", err)
	}

	if err := sendMessageSNAC(msgCh, chatCtx.cookie, chatMsg.ScreenName, botResponse); err != nil {
		return fmt.Errorf("unable to send response: %w", err)
	}

	// store this exchange to be used as context for next bot request
	chatCtx.lastExchange[0] = userMessage
	chatCtx.lastExchange[1] = botResponse

	if chatCtx.warnCount == 3 {
		sendWarningSNAC(msgCh, chatMsg.ScreenName)
	}

	return nil
}

// exchangeMessages receives an IM and responds with a bot message.
func exchangeMessages(
	logger *slog.Logger,
	msgCh chan wire.SNACMessage,
	flapBody *bytes.Buffer,
	chatContexts map[string]*chatContext,
	chatBot ChatBot,
	config config.Config,
) error {

	msgSNAC := wire.SNAC_0x04_0x07_ICBMChannelMsgToClient{}
	if err := wire.Unmarshal(&msgSNAC, flapBody); err != nil {
		return err
	}

	if _, ok := chatContexts[msgSNAC.ScreenName]; !ok {
		// this is the first message received from this user
		chatContexts[msgSNAC.ScreenName] = &chatContext{
			cookie:       msgSNAC.Cookie,
			semaphore:    make(chan struct{}, 1),
			lastExchange: new([2]string),
			limiter:      rate.NewLimiter(rate.Every(time.Minute), config.MaxMsgPerMin),
		}
	}

	// Retrieve chat context for current user.
	chatCtx := chatContexts[msgSNAC.ScreenName]

	// Update context with the latest conversation unique ID.
	chatCtx.cookie = msgSNAC.Cookie

	// Ignore this message if the bot is currently processing a message
	// exchange.
	if !chatCtx.tryLock() {
		return nil // currently responding to user, drop message
	}

	var messageSent bool
	defer func() {
		// Ensure the lock is released if this function exits prematurely,
		// before the message send goroutine, which normally releases the lock
		if !messageSent {
			chatCtx.releaseLock()
		}
	}()

	// Ensure the user is not chatting with the bot too quickly. When they
	// reach the rate limit threshold, inform the user that they are sending
	// messages too quickly and ignore subsequent messages until the rate limit
	// window passes.
	if hitRateLimit := enforceRateLimit(logger, msgCh, chatCtx, msgSNAC); hitRateLimit {
		logger.Info("user hit message rate limit", "screen_name", msgSNAC.ScreenName)
		return nil
	}

	// Get the message text buried in the SNAC payload.
	if _, hasIMData := msgSNAC.TLVRestBlock.Slice(wire.ICBMTLVAOLIMData); !hasIMData {
		logger.Debug("received ICBMChannelMsgToClient with no AOLIMData")
		return nil
	}
	msgText, err := msgSNAC.ExtractMessageText()
	if err != nil {
		return err
	}

	// Strip HTML formatting so that we don't confuse the bot.
	msgText = stripHTMLTags(msgText)

	// Make sure the message is not too big in order to minimize cost. OpenAI
	// charges per token (which is effectively a word).
	if hitMsgSizeLimit := enforceMsgSizeLimit(logger, msgText, msgCh, msgSNAC, config); hitMsgSizeLimit {
		logger.Info("user hit message size limit", "screen_name", msgSNAC.ScreenName)
		return nil
	}

	messageSent = true

	go func() {
		defer chatCtx.releaseLock()

		if _, wantsEvents := msgSNAC.TLVRestBlock.Slice(wire.ICBMTLVWantEvents); wantsEvents {
			// Tell the client that the bot is "typing". Provides a visual
			// indicator in the IM window that something is happening.
			sendTypingEventSNAC(msgSNAC, msgCh)
		}

		// Get the bot's response to this message.
		botResponse, err := chatBot.ExchangeMessage(msgText, *chatCtx.lastExchange)
		if err != nil {
			logger.Error("unable to get response from bot", "err", err.Error())
			return
		}

		// Send the bot's response.
		if err := sendMessageSNAC(msgCh, msgSNAC.Cookie, msgSNAC.ScreenName, botResponse); err != nil {
			logger.Error("unable to send response", "err", err.Error())
			return
		}

		// Save this interaction for use as context in the next bot request.
		chatCtx.lastExchange[0] = msgText
		chatCtx.lastExchange[1] = botResponse

		logger.Info("message exchange", "screen_name", msgSNAC.ScreenName, "incoming", msgText, "outgoing", botResponse)
	}()

	return nil
}

func enforceMsgSizeLimit(
	logger *slog.Logger,
	text string,
	msgCh chan wire.SNACMessage,
	msgSNAC wire.SNAC_0x04_0x07_ICBMChannelMsgToClient,
	config config.Config,
) bool {

	tooLong := exceedsMsgSizeLimit(text, config)
	if tooLong {
		botResponse := "Your message is too long for me! I am but a simple bot!"
		if err := sendMessageSNAC(msgCh, msgSNAC.Cookie, msgSNAC.ScreenName, botResponse); err != nil {
			logger.Error("unable to send size limit warning", "err", err.Error())
		}
	}
	return tooLong
}

func exceedsMsgSizeLimit(text string, config config.Config) bool {
	words := strings.Fields(text)
	if len(words) > config.WordCountLimit {
		return true
	}
	for _, word := range words {
		if len(word) > config.WordLengthLimit {
			return true
		}
	}
	return false
}

func enforceRateLimit(
	logger *slog.Logger,
	msgCh chan wire.SNACMessage,
	chatCtx *chatContext,
	msgSNAC wire.SNAC_0x04_0x07_ICBMChannelMsgToClient,
) bool {
	if !chatCtx.limiter.Allow() {
		if !chatCtx.rateLimited {
			chatCtx.rateLimited = true
			go func() {
				botResponse := "You're sending me too many messages! Slow down!"
				if err := sendMessageSNAC(msgCh, msgSNAC.Cookie, msgSNAC.ScreenName, botResponse); err != nil {
					logger.Error("unable to send rate limit limit warning", "err", err.Error())
					return
				}
			}()
		}
		return true
	}
	// clear the rate limit flag that may have been set previously
	chatCtx.rateLimited = false
	return false
}

func sendMessageSNAC(msgCh chan<- wire.SNACMessage, cookie uint64, screenName string, response string) error {
	msgFrame := wire.SNACFrame{
		FoodGroup: wire.ICBM,
		SubGroup:  wire.ICBMChannelMsgToHost,
	}
	responseSNAC := wire.SNAC_0x04_0x06_ICBMChannelMsgToHost{
		Cookie:     cookie,
		ChannelID:  1,
		ScreenName: screenName,
	}

	response = fmt.Sprintf(`<HTML><BODY BGCOLOR="#ffffff">%s</BODY></HTML>`, response)
	if err := responseSNAC.ComposeMessage(response); err != nil {
		return fmt.Errorf("unable to compose message: %w", err)
	}

	msgCh <- wire.SNACMessage{
		Frame: msgFrame,
		Body:  responseSNAC,
	}

	return nil
}

func sendWarningSNAC(msgCh chan<- wire.SNACMessage, screenName string) {
	msgCh <- wire.SNACMessage{
		Frame: wire.SNACFrame{
			FoodGroup: wire.ICBM,
			SubGroup:  wire.ICBMEvilRequest,
		},
		Body: wire.SNAC_0x04_0x08_ICBMEvilRequest{
			ScreenName: screenName,
		},
	}
}

func sendTypingEventSNAC(chatMsg wire.SNAC_0x04_0x07_ICBMChannelMsgToClient, msgCh chan<- wire.SNACMessage) {
	msgCh <- wire.SNACMessage{
		Frame: wire.SNACFrame{
			FoodGroup: wire.ICBM,
			SubGroup:  wire.ICBMClientEvent,
		},
		Body: wire.SNAC_0x04_0x14_ICBMClientEvent{
			Cookie:     chatMsg.Cookie,
			ChannelID:  chatMsg.ChannelID,
			ScreenName: chatMsg.ScreenName,
			Event:      0x0002, // indicates "typing begun"
		},
	}
}

var stripHTMLRegex = regexp.MustCompile("<[^>]*>")

func stripHTMLTags(input string) string {
	return stripHTMLRegex.ReplaceAllString(input, "")
}
