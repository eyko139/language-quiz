package words

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eyko139-language-app/cmd/env"

	core "github.com/Azure/azure-sdk-for-go/sdk/azcore"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
)

type GPTInt interface {
	GetTranslations(word string) []GptWord
}

type GPT struct {
	Client *azopenai.Client
}

func NewGPT(env *env.Env) (*GPT, error) {
	client, err := initClient(env)
	if err != nil {
		return nil, err
	}
	return &GPT{
		Client: client,
	}, nil
}

var client *azopenai.Client
var messages []azopenai.ChatRequestMessageClassification

func initClient(env *env.Env) (*azopenai.Client, error) {
	messages = []azopenai.ChatRequestMessageClassification{
		// You set the tone and rules of the conversation with a prompt as the system role.

		// The user asks a question
		&azopenai.ChatRequestSystemMessage{Content: azopenai.NewChatRequestSystemMessageContent(env.GPT_PROMPT)},
	}

	modelDeploymentID := "gpt-4o"
	maxTokens := int32(400)
	keyCredential := core.NewKeyCredential(env.AZ_OPENAI_KEY)

	client, _ = azopenai.NewClientWithKeyCredential(env.AZ_OPENAI_ENDPOINT, keyCredential, nil)

	_, _ = client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		// This is a conversation in progress.
		Messages:       messages,
		DeploymentName: &modelDeploymentID,
		MaxTokens:      &maxTokens,
	}, nil)

	return client, nil
}

type GptWord struct {
	Word        string `json:"word"`
	T_1         string `json:"t_1"`
	T_2         string `json:"t_2"`
	T_3         string `json:"t_3"`
	T_4         string `json:"t_4"`
	Description string `json:"description"`
}

func (g *GPT) GetTranslations(words []string) []GptWord {

	modelDeploymentID := "gpt-4o"
	maxTokens := int32(2000)

	var translations []GptWord

	messages = append(messages, &azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(strings.Join(words[:], ","))})

	resp, _ := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		// This is a conversation in progress.
		// NOTE: all messages count against token usage for this API.
		Messages:       messages,
		DeploymentName: &modelDeploymentID,
		MaxTokens:      &maxTokens,
	}, nil)

	err := json.Unmarshal([]byte(*resp.Choices[0].Message.Content), &translations)
	if err != nil {
		fmt.Println(err)
	}
	return translations
}

func (g *GPT) TranslateText(text string) []GptWord {

	modelDeploymentID := "gpt-4o"
	maxTokens := int32(5000)

	var translations []GptWord

	messages = append(messages, &azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(text)})

	resp, _ := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		// This is a conversation in progress.
		// NOTE: all messages count against token usage for this API.
		Messages:       messages,
		DeploymentName: &modelDeploymentID,
		MaxTokens:      &maxTokens,
	}, nil)

	json.Unmarshal([]byte(*resp.Choices[0].Message.Content), &translations)
	return translations
}
