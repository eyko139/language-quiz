package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (a *App) Routes() http.Handler {

	router := httprouter.New()

	router.Handler(http.MethodGet, "/all-words", a.getHome())
	router.Handler(http.MethodPost, "/save-word", a.wordPost())
	router.Handler(http.MethodGet, "/translate", a.translate())
	router.Handler(http.MethodPost, "/bot-message", a.botHook())

	standard := alice.New(a.CorsRequest, a.LogRequests)

	return standard.Then(router)

}
