package main

import (
	"net/http"
	"time"
    "fmt"

	"github.com/eyko139-language-app/cmd/env"
)


func main() {

	appEnv := env.New()
    
    fmt.Println("setting up...")
	app := NewApp(appEnv)


	srv := http.Server{
		ErrorLog:     app.Errlog,
		Addr:         ":" + appEnv.AppPort,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

    app.Infolog.Printf("Starting server on %s", srv.Addr)
    err := srv.ListenAndServe()
	app.Errlog.Fatal(err)

}


