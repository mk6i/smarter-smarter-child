package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/mk6i/retro-aim-server/wire"

	"github.com/mk6i/smarter-smarter-child/bot"
	"github.com/mk6i/smarter-smarter-child/client"
	"github.com/mk6i/smarter-smarter-child/config"
)

func main() {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to process app config: %s\n", err.Error())
		os.Exit(1)
	}

	logger := NewLogger(cfg)

	bosHost, authCookie, err := func() (string, string, error) {
		host := net.JoinHostPort(cfg.OSCARHost, cfg.OSCARPort)
		conn, err := net.Dial("tcp", host)
		if err != nil {
			return "", "", fmt.Errorf("unable to dial into auth host: %w", err)
		}
		defer func() {
			logger.Debug("disconnected from auth service", "host", host)
			conn.Close()
		}()

		logger.Debug("connected to auth service", "host", host)

		flapc := wire.NewFlapClient(0, conn, conn)
		host, authCookie, err := client.Authenticate(flapc, cfg.ScreenName, cfg.Password)
		if err == nil {
			logger.Debug("authentication succeeded, proceeding to BOS host", "host", host, "authCookie", authCookie)
		}
		return host, authCookie, err
	}()

	if err != nil {
		logger.Error("authentication failed", "err", err.Error())
		os.Exit(1)
	}

	err = func() error {
		conn, err := net.Dial("tcp", bosHost)
		if err != nil {
			return err
		}
		defer conn.Close()

		logger.Info("connected to BOS server", "host", bosHost)

		flapc := wire.NewFlapClient(0, conn, conn)

		var chatBot client.ChatBot
		if cfg.OfflineMode {
			logger.Debug("offline mode enabled, using local chatbot backend")
			chatBot = bot.NewStaticChatBot()
		} else {
			logger.Debug("using ChatGPT chatbot backend")
			chatBot = bot.NewChatGPTBot(cfg.OpenAIKey)
		}
		return client.Chat(logger, flapc, authCookie, chatBot, cfg)
	}()

	if err != nil {
		logger.Error("chat failed", "err", err.Error())
		os.Exit(1)
	}
}

func NewLogger(cfg config.Config) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	case "info":
		fallthrough
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}
	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}
