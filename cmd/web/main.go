package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eyko139-language-app/cmd/env"
)

func main() {

	appEnv := env.New()

	fmt.Println("setting up...")
	app := NewApp(appEnv)

	srv := http.Server{
		ErrorLog:     app.ErrorLog,
		Addr:         ":" + appEnv.AppPort,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.InfoLog.Printf("Starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	app.ErrorLog.Fatal(err)

}
