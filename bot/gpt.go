package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/mk6i/smarter-smarter-child/config"
)

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature"`
	TopP        float64   `json:"top_p"`
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

type errorResponse struct {
	Error struct {
		Message string      `json:"message"`
		Type    string      `json:"type"`
		Param   interface{} `json:"param"`
		Code    string      `json:"code"`
	} `json:"error"`
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

func NewChatGPTBot(cfg config.Config) *ChatGPTChatBot {
	return &ChatGPTChatBot{
		secretKey:   cfg.OpenAIKey,
		prompt:      cfg.BotPrompt,
		model:       cfg.Model,
		temperature: cfg.Temperature,
		topP:        cfg.TopP,
		apiURL:      cfg.APIUrl,
		client:      &http.Client{},
		r:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type ChatGPTChatBot struct {
	r           *rand.Rand
	secretKey   string
	prompt      string
	model       string
	temperature float64
	topP        float64
	apiURL      string
	client      *http.Client
}

func (g *ChatGPTChatBot) ExchangeMessage(send string, lastExchange [2]string) (receive string, err error) {
	messages := []message{
		{
			Role:    "system",
			Content: g.prompt,
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
		Model:       g.model,
		Messages:    messages,
		Temperature: g.temperature,
		TopP:        g.topP,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", g.apiURL, bytes.NewBuffer(jsonData))
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

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
		case http.StatusForbidden:
		case http.StatusTooManyRequests:
		case http.StatusNotFound:
			var response errorResponse
			if err := json.Unmarshal(body, &response); err != nil {
				return "", err
			}
			return "", fmt.Errorf("error from upstream api: %s", response.Error.Message)
		default:
			return "", fmt.Errorf("unknown error from upstream api: %d", resp.StatusCode)
		}
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
