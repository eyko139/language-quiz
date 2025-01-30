package env

import (
	"github.com/spf13/viper"
)

type Env struct {
	AppPort            string
	DB_FILE            string
	GPT_PROMPT         string
	AZ_OPENAI_KEY      string
	AZ_OPENAI_ENDPOINT string
}

func New() *Env {

	viper.BindEnv("DB_FILE")
	viper.SetDefault("DB_FILE", "/home/luk/projects/priv/language_app/words.db")

    viper.BindEnv("GPT_PROMPT")
    viper.SetDefault("GPT_PROMPT", "You return 4 alternives of english translations to a german word. ONly one of them should be correct. Answer in json format. Do not start with ```. If i provide multiple words, return an array. The format of each translation should be like this {word: the initial word i gave you, t_1: the correct translation, t_2: one of the incorrect translations, t_3: one of the incorrect translations, t_4: one of the incorrect translations}")

    viper.BindEnv("AZ_OPENAI_KEY")

    viper.BindEnv("AZ_OPENAI_ENDPOINT")

	return &Env{
		AppPort: "8080",
        DB_FILE: viper.GetString("DB_FILE"),
        GPT_PROMPT: viper.GetString("GPT_PROMPT"),
        AZ_OPENAI_KEY: viper.GetString("AZ_OPENAI_KEY"),
        AZ_OPENAI_ENDPOINT: viper.GetString("AZ_OPENAI_ENDPOINT"),
	}
}
