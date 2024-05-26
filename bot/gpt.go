package bot

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type completionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Usage   usage    `json:"usage"`
	Choices []choice `json:"choices"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type choice struct {
	Message      message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

func NewChatGPTBot(secretKey string) *ChatGPTChatBot {
	return &ChatGPTChatBot{
		secretKey: secretKey,
		client:    &http.Client{},
		r:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type ChatGPTChatBot struct {
	r         *rand.Rand
	secretKey string
	client    *http.Client
}

func (g *ChatGPTChatBot) ExchangeMessage(send string, lastExchange [2]string) (receive string, err error) {
	messages := []message{
		{
			Role:    "system",
			Content: "You are SmarterChild, a dumb AIM chatbot.",
		},
	}
	if lastExchange[0] != "" {
		messages = append(messages, message{
			Role:    "user",
			Content: lastExchange[0],
		})
		messages = append(messages, message{
			Role:    "assistant",
			Content: lastExchange[1],
		})
	}
	messages = append(messages, message{
		Role:    "user",
		Content: send,
	})

	data := chatRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    messages,
		Temperature: 0.7,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.secretKey)

	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response completionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "No response available.", nil
}
