package setup

import (
	"github.com/melody-mood/config"
	openai "github.com/sashabaranov/go-openai"
)

func InitOpenAIService() *openai.Client {
	conf := config.GetConfig()
	client := openai.NewClient(conf.OpenAI.APIKey)
	return client
}
