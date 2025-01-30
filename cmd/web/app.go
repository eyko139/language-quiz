package main

import (
	"log"
	"os"

	"github.com/eyko139-language-app/cmd/env"
	"github.com/eyko139-language-app/gpt"
)

type App struct {
	Errlog    *log.Logger
	Infolog   *log.Logger
	WordModel words.WordModelInt
}

func NewApp(env *env.Env) *App {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	wordModel, err := words.NewWordModel(env)

	if err != nil {
		errLog.Fatal(err)
	}

	return &App{
		Errlog:    errLog,
		Infolog:   infoLog,
		WordModel: wordModel,
	}
}
