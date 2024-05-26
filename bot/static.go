package bot

import (
	"math/rand"
	"time"
)

func NewStaticChatBot() *StaticChatBot {
	return &StaticChatBot{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type StaticChatBot struct {
	r *rand.Rand
}

func (c *StaticChatBot) ExchangeMessage(send string, exchange [2]string) (receive string, err error) {
	time.Sleep(time.Duration(c.r.Intn(1000)) * time.Millisecond)
	responses := []string{
		"hi2u",
		"a/s/l?",
		"brb",
		"lol",
		"rofl",
		"ttyl",
		"omg",
		"g2g",
		"idk",
		"bbl",
	}
	receive = responses[c.r.Intn(len(responses))]
	return
}
