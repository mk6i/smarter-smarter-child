package client

import (
	"errors"
	"fmt"

	"github.com/mk6i/retro-aim-server/wire"
)

// Authenticate performs the BUCP auth flow with the OSCAR auth server. Upon
// successful login, it returns a host name and auth cookie for connecting to
// and authenticating with the BOS service.
func Authenticate(flapc FlapClient, screenName string, password string) (string, string, error) {
	if _, err := flapc.ReceiveSignonFrame(); err != nil {
		return "", "", fmt.Errorf("unable to receive signon frame: %w", err)
	}
	if err := flapc.SendSignonFrame(nil); err != nil {
		return "", "", fmt.Errorf("unable to send signon frame: %w", err)
	}

	challengeRequest := wire.SNAC_0x17_0x06_BUCPChallengeRequest{}
	challengeRequest.Append(wire.NewTLV(wire.LoginTLVTagsScreenName, screenName))
	if err := flapc.SendSNAC(wire.SNACFrame{
		FoodGroup: wire.BUCP,
		SubGroup:  wire.BUCPChallengeRequest,
	}, challengeRequest); err != nil {
		return "", "", fmt.Errorf("unable to send SNAC(0x17,0x06): %w", err)
	}

	challengeResponse := &wire.SNAC_0x17_0x07_BUCPChallengeResponse{}
	if err := flapc.ReceiveSNAC(&wire.SNACFrame{}, challengeResponse); err != nil {
		return "", "", fmt.Errorf("unable to receive SNAC(0x17,0x07): %w", err)
	}

	loginRequest := wire.SNAC_0x17_0x02_BUCPLoginRequest{}
	loginRequest.Append(wire.NewTLV(wire.LoginTLVTagsScreenName, screenName))
	loginRequest.Append(wire.NewTLV(wire.LoginTLVTagsPasswordHash,
		wire.StrongMD5PasswordHash(password, challengeResponse.AuthKey)))
	if err := flapc.SendSNAC(wire.SNACFrame{
		FoodGroup: wire.BUCP,
		SubGroup:  wire.BUCPLoginRequest,
	}, loginRequest); err != nil {
		return "", "", fmt.Errorf("unable to send SNAC(0x17,0x02): %w", err)
	}

	loginRespSNAC := wire.SNAC_0x17_0x03_BUCPLoginResponse{}
	if err := flapc.ReceiveSNAC(&wire.SNACFrame{}, &loginRespSNAC); err != nil {
		return "", "", fmt.Errorf("unable to receive SNAC(0x17,0x03): %w", err)
	}

	if code, hasErr := loginRespSNAC.Uint16(wire.LoginTLVTagsErrorSubcode); hasErr {
		switch code {
		case wire.LoginErrInvalidUsernameOrPassword:
			return "", "", fmt.Errorf("error code from SNAC(0x17,0x03): invalid username or password")
		default:
			return "", "", fmt.Errorf("error code from SNAC(0x17,0x03): %d", code)
		}
	}

	host, hasHostname := loginRespSNAC.String(wire.LoginTLVTagsReconnectHere)
	if !hasHostname {
		return "", "", errors.New("SNAC(0x17,0x03) does not contain a hostname TLV")
	}

	authCookie, hasAuthCookie := loginRespSNAC.String(wire.LoginTLVTagsAuthorizationCookie)
	if !hasAuthCookie {
		return "", "", errors.New("SNAC(0x17,0x03) does not contain an auth cookie TLV")
	}

	return host, authCookie, nil
}
