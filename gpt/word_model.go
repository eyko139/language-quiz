package words

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/eyko139-language-app/cmd/env"
	_ "github.com/mattn/go-sqlite3"
)

const (
	createDB = `
        CREATE TABLE IF NOT EXISTS words (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        time DATETIME NOT NULL,
        word TEXT,
        t_1 TEXT,
        t_2 TEXT,
        t_3 TEXT,
        t_4 TEXT
        );
        `
)

type Word struct {
    ID int `json:"id" db:"id"`
	Word         string    `json:"word" db:"word"`
	Translation1 string    `json:"t_1" db:"t_1"`
	Translation2 string    `json:"t_2" db:"t_2"`
	Translation3 string    `json:"t_3" db:"t_3"`
	Translation4 string    `json:"t_4" db:"t_4"`
	Time         time.Time `json:"time" db:"time"`
}

type WordModelInt interface {
	GetWord(word string) (Word, error)
	GetAllWords() ([]Word, error)
	AddWord(word string) (int64, error)
	Translate() ([]string, error)
}

type WordModel struct {
	Db  *sql.DB
	env *env.Env
	GPT *GPT
}

func NewWordModel(env *env.Env) (*WordModel, error) {

    fmt.Println("init db")
	db, err := sql.Open("sqlite3", env.DB_FILE)

	if err != nil {
		return nil, err
	}

    fmt.Println("creating table")
	_, err = db.Exec(createDB)


	if err != nil {
		return nil, err
	}

    fmt.Println("init GPT client")
	gpt, err := NewGPT(env)

	if err != nil {
		return nil, err
	}
    fmt.Println("GPT client finished")

	return &WordModel{Db: db, env: env, GPT: gpt}, nil
}

func (m *WordModel) GetWord(word string) (Word, error) {
	return Word{}, nil
}

func (m *WordModel) GetAllWords() ([]Word, error) {
	rows, err := m.Db.Query("select * from words where t_1 is not null;")

	if err != nil {
		return nil, err
	}

	var words []Word

	for rows.Next() {
		var word Word
        err := rows.Scan(&word.ID, &word.Time, &word.Word, &word.Translation1, &word.Translation2, &word.Translation3, &word.Translation4,)
        if err != nil {
            return nil, err
        }
        words = append(words, word)
	}

	return words, nil
}

func (m *WordModel) AddWord(newWord string) (int64, error) {

	insertWord := "INSERT INTO words (time, word, t_1, t_2, t_3, t_4) VALUES (?, ?, ?, ?, ?, ?)"

	word := Word{
		Time: time.Now(),
		Word: newWord,
	}

	res, err := m.Db.Exec(insertWord, time.Now(), word.Word, nil, nil, nil, nil)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *WordModel) Translate() ([]string, error) {
	rows, err := m.Db.Query("select id, time, word from words where t_1 is null;")

	if err != nil {
		return nil, err
	}

	wordsToTranslate := []string{}

	for rows.Next() {
		var id int
		var word, time string

		err := rows.Scan(&id, &time, &word)
		if err != nil {
			log.Fatalf("Failed to scan rows %v", err)
		}
		wordsToTranslate = append(wordsToTranslate, word)
	}

	defer rows.Close()

	translations := m.GPT.GetTranslations(wordsToTranslate)

	updateWord := "update words set t_1 = ?, t_2 = ?, t_3 = ?, t_4 = ? where word = ?;"

	for _, translation := range translations {
		fmt.Printf("%v", translation)
		res, err := m.Db.Exec(updateWord, translation.T_1, translation.T_2, translation.T_3, translation.T_4, translation.Word)
		fmt.Println(res)
		if err != nil {
			return nil, err
		}
	}

	rows, err = m.Db.Query("select id, time, word from words where t_1 is not null;")

	if err != nil {
		return nil, err
	}

	translatedWords := []string{}

	for rows.Next() {
		var id int
		var word, time string

		err := rows.Scan(&id, &time, &word)
		if err != nil {
			log.Fatalf("Failed to scan rows %v", err)
		}
		translatedWords = append(translatedWords, word)
	}
	return translatedWords, nil
}
