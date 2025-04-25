package pkg

import (
	"context"
	"encoding/json"

	"github.com/melody-mood/constants"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Define the structure for the OpenAI response
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func CreateOpenAIMessage(ctx context.Context, client *openai.Client, prompt string) (resp string, err error) {
	req := openai.ChatCompletionRequest{
		Model: constants.OpenAIModel, // or "gpt-3.5-turbo" if you prefer
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	res, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	responseContent := res.Choices[0].Message.Content

	logrus.Printf("[OpenAI Prompt] %s", prompt)
	logrus.Printf("[OpenAI Response] %s", responseContent)

	return responseContent, nil
}

func ParseToStruct[T any](jsonString string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
