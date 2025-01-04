package seeder

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
	client *openai.Client
}

func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{
		client: openai.NewClient(apiKey),
	}
}

func (ai *OpenAIService) queryllm(systemPrompt string, userPrompt string) (string, error) {

	res, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("error generating message openai llm: %w", err)
	}

	response := strings.TrimPrefix(res.Choices[0].Message.Content, "```json\n")
	response = strings.TrimSuffix(response, "```")

	return response, nil
}
