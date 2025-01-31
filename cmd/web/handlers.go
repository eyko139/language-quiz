package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/eyko139-language-app/models"
	"github.com/julienschmidt/httprouter"

	mod "github.com/eyko139-language-app/gpt"
)

func (a *App) getHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		translatedWords, totalWordsAmount, percentageTranslated, err := a.WordModel.GetAllWords()
		if err != nil {
			a.ErrorLog.Println(err)
			_, err = w.Write([]byte("error getting all translatedWords"))
			if err != nil {
				a.ErrorLog.Println(err)
			}
		}

		w.Header().Add("Content-Type", "application/json")
		enc := json.NewEncoder(w)

		homeView := struct {
			TranslatedWords        []mod.Word `json:"allWords"`
			PercentageUntranslated float64    `json:"percentageUntranslated"`
			TotalWords             int        `json:"totalWords"`
		}{
			TranslatedWords:        translatedWords,
			PercentageUntranslated: percentageTranslated,
			TotalWords:             totalWordsAmount,
		}

		err = enc.Encode(homeView)
		if err != nil {
			a.ErrorLog.Println(err)
		}
	}
}

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
			a.ErrorLog.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte("Error adding word"))
			if err != nil {
				a.ErrorLog.Println(err)
			}
		}

		_, err = w.Write([]byte(myAns.Message))
		if err != nil {
			a.ErrorLog.Println(err)
		}
	}
}

func (a *App) translate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		translations, err := a.WordModel.Translate()
		if err != nil {
			a.ErrorLog.Println("error during translation")
			_, err = w.Write([]byte("error during translation"))
			if err != nil {
				a.ErrorLog.Println(err)
			}
		}

		enc := json.NewEncoder(w)

		err = enc.Encode(translations)
		if err != nil {
			a.ErrorLog.Println(err)
		}
	}
}

func (a *App) botHook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var botResponse models.BotResponse

		if err := json.NewDecoder(r.Body).Decode(&botResponse); err != nil {
			a.ErrorLog.Printf("Error parsing bot response: %s", err)
		}
		if botResponse.Message.From.IsBot {
			return
		}

		err := a.WordModel.AddText(botResponse.Message.Text)
		if err != nil {
			a.ErrorLog.Println(err)
		}
	}
}
func (a *App) submitAnswerPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
        id, err := strconv.Atoi(params.ByName("id"))
        if err != nil {
            a.ErrorLog.Printf("Failed to parse int wordID %s", err)
        }

        translation := params.ByName("translation")
        a.InfoLog.Printf("submited answer for %s: %s", id, translation)
        err = a.WordModel.EvalAnswer(translation, id)
        if err != nil {
            a.ErrorLog.Printf("Error submitting answer %s", err)
        }
	}
}
