package main

import (
	"log"
	"os"

	"github.com/eyko139-language-app/cmd/env"
	"github.com/eyko139-language-app/gpt"
	"github.com/gorilla/websocket"
)

type App struct {
	ErrorLog  *log.Logger
	InfoLog   *log.Logger
	WordModel words.WordModelInt
	WsConn    *websocket.Conn
}

func NewApp(env *env.Env) *App {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	wordModel, err := words.NewWordModel(env)

	if err != nil {
		errLog.Fatal(err)
	}

	return &App{
		ErrorLog:  errLog,
		InfoLog:   infoLog,
		WordModel: wordModel,
	}
}
