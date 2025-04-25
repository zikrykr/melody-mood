package setup

import (
	"github.com/melody-mood/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func InitOpenAIService() *openai.Client {
	conf := config.GetConfig()
	client := openai.NewClient(option.WithAPIKey(conf.OpenAI.APIKey))
	return &client
}
