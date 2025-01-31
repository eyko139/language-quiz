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
        word TEXT unique,
        t_1 TEXT NULL,
        t_2 TEXT NULL,
        t_3 TEXT NULL,
        t_4 TEXT NULL,
        description TEXT NULL,
        timesQuizzed INTEGER DEFAULT 0,
        timesCorrect INTEGER DEFAULT 0, 
        timesFalse INTEGER DEFAULT 0
        );
        `
)

type Word struct {
	ID           int       `json:"id" db:"id"`
	Word         string    `json:"word" db:"word"`
	Translation1 string    `json:"t_1" db:"t_1"`
	Translation2 string    `json:"t_2" db:"t_2"`
	Translation3 string    `json:"t_3" db:"t_3"`
	Translation4 string    `json:"t_4" db:"t_4"`
	Time         time.Time `json:"time" db:"time"`
	Description  string    `json:"description" db:"description"`
	TimesQuizzed int       `json:"timesQuizzed" db:"timesQuizzed"`
	TimesCorrect int       `json:"timesCorrect" db:"timesCorrect"`
	TimesFalse   int       `json:"timesFalse" db:"timesFalse"`
}

type WordModelInt interface {
	GetWord(word string) (Word, error)
	GetAllWords() ([]Word, int, float64, error)
	AddWord(word string) (int64, error)
	Translate() ([]string, error)
	AddText(text string) error
	EvalAnswer(selectedTranslation string, wordId int) error
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

func (m *WordModel) GetAllWords() ([]Word, int, float64, error) {
	rows, err := m.Db.Query("select * from words where t_1 is not null;")

	var translated, untranslated sql.NullInt32
	var percentage sql.NullFloat64

	err = m.Db.QueryRow(`
    SELECT 
        COUNT(CASE WHEN t_1 IS NOT NULL THEN 1 END) AS translated,
        COUNT(CASE WHEN t_1 IS NULL THEN 1 END) AS untranslated,
        (COUNT(CASE WHEN t_1 IS NULL THEN 1 END) * 100.0) / COUNT(*) AS percentage_untranslated
    FROM words;
`).Scan(&translated, &untranslated, &percentage)

	if err != nil {
		return nil, 0, 0, err
	}

	var words []Word

	for rows.Next() {
		var id int
		var wordString string
		var time time.Time
		var t1, t2, t3, t4, description sql.NullString
		var timesQuizzed, timesCorrect, timesFalse int
		err := rows.Scan(&id, &time, &wordString, &t1, &t2, &t3, &t4, &description, &timesQuizzed, &timesCorrect, &timesFalse)
		if err != nil {
			return nil, 0, 0, err
		}
		var word = Word{
			ID:           id,
			Word:         wordString,
			Translation1: t1.String,
			Translation2: t2.String,
			Translation3: t3.String,
			Translation4: t4.String,
			Time:         time,
			Description:  description.String,
			TimesQuizzed: timesQuizzed,
			TimesCorrect: timesCorrect,
			TimesFalse:   timesFalse,
		}
		words = append(words, word)
	}

	if translated.Valid && untranslated.Valid && percentage.Valid {
		return words, int(translated.Int32 + untranslated.Int32), float64(percentage.Float64), nil
	}

	return words, 0, 0, nil
}

func (m *WordModel) AddText(text string) error {
	translations := m.GPT.TranslateText(text)

	return m.saveWords(translations)
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

	fmt.Println(wordsToTranslate)
	translations := m.GPT.GetTranslations(wordsToTranslate)

	return m.updateWord(translations)

}

func (m *WordModel) saveWords(translation []GptWord) error {
	insertWord := "INSERT INTO words (time, word, t_1, t_2, t_3, t_4, description) VALUES (?, ?, ?, ?, ?, ?, ?)"
	for _, t := range translation {
		_, err := m.Db.Exec(insertWord, time.Now(), t.Word, t.T_1, t.T_2, t.T_3, t.T_4, t.Description)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *WordModel) updateWord(translations []GptWord) ([]string, error) {

	updateWord := "update words set t_1 = ?, t_2 = ?, t_3 = ?, t_4 = ? where word = ?;"

	for _, translation := range translations {
		fmt.Printf("%v", translation)
		res, err := m.Db.Exec(updateWord, translation.T_1, translation.T_2, translation.T_3, translation.T_4, translation.Word)
		fmt.Println(res)
		if err != nil {
			return nil, err
		}
	}

	rows, err := m.Db.Query("select id, time, word from words where t_1 is not null;")

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

func (m *WordModel) EvalAnswer(selectedTranslation string, wordId int) error {
	updateWord := "update words set timesQuizzed = ?, timesCorrect = ?, timesFalse = ? where id = ?;"

	var word Word

	if err := m.Db.QueryRow("select timesQuizzed, timesCorrect, timesFalse, t_1 from words where id = ?", wordId).Scan(&word.TimesQuizzed, &word.TimesCorrect, &word.TimesFalse, &word.Translation1); err != nil {
		return err
	}

	fmt.Printf("%+v", word)

	if word.Translation1 == selectedTranslation {
		_, err := m.Db.Exec(updateWord, word.TimesQuizzed+1, word.TimesCorrect+1, word.TimesFalse, wordId)
		if err != nil {
			return err
		}
	} else {
		_, err := m.Db.Exec(updateWord, word.TimesQuizzed+1, word.TimesCorrect, word.TimesFalse+1, wordId)
		if err != nil {
			return err
		}
	}

	return nil
}
