package client

import (
	"github.com/mk6i/retro-aim-server/wire"
)

type ChatBot interface {
	ExchangeMessage(send string, exchange [2]string) (receive string, err error)
}

type FlapClient interface {
	ReceiveFLAP() (frame wire.FLAPFrame, err error)
	ReceiveSNAC(frame *wire.SNACFrame, body any) error
	ReceiveSignonFrame() (wire.FLAPSignonFrame, error)
	SendSNAC(frame wire.SNACFrame, body any) error
	SendSignonFrame(tlvs []wire.TLV) error
}
