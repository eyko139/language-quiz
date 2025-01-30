package main

import (
	"encoding/json"
    "io"
	"net/http"
)

func (a *App) allWords() http.HandlerFunc { 
    return func(w http.ResponseWriter, r *http.Request) {
        words, err := a.WordModel.GetAllWords()
        if err != nil {
            a.Errlog.Println(err)
            w.Write([]byte("error getting all words"))
        }

        enc := json.NewEncoder(w)

        enc.Encode(words)
    }}

func (a *App) wordPost() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}

			type answer struct {
				Message string `json:"message"`
			}

			var myAns answer
			// Unmarshal the request body
			err = json.Unmarshal(body, &myAns)

            _, err = a.WordModel.AddWord(myAns.Message)

            if err != nil {
                a.Errlog.Println(err.Error())
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte("Error adding word"))
            }

			w.Write([]byte(myAns.Message))
    }
}

func (a *App) translate() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    translations, err := a.WordModel.Translate()
    if err != nil {
        a.Errlog.Println("error during translation")
        w.Write([]byte("error during translation"))
    }

    enc := json.NewEncoder(w)

    enc.Encode(translations)
}
}
