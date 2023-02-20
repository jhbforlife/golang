package translate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// Path to supported languages database
var dbPath = "./languages.db"

// Error returned if a language is either invalid or not supported.
var ErrInvalidLang = errors.New("invalid or unsupported language: ")

// Error returned if no translation language is provided.
var ErrNoToLang = errors.New("no translation language provided")

// Error returned if no translation text is provided.
var ErrNoText = errors.New("no translation text provided")

// Returns a slice of supported language names.
func SupportedLanguages() ([]string, error) {
	return getSupportedLanguages()
}

// Translate text to another language. Returns an error if either
// to or text are empty, or if either language is invalid or unsupported.
// If from is not provided, the source language will be assumed. Returns a slice of
// possible translations, ranked from hightest to lowest confidence.
func TranslateText(from, to, text string) ([]string, error) {

	// Check if to or text are empty before sending a translation request
	if isEmptyString(to) {
		return nil, ErrNoToLang
	}
	if isEmptyString(text) {
		return nil, ErrNoText
	}

	// Initialize Cloud Translation client
	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Verify database exists and is up to date
	if verifyDB(); err != nil {
		return nil, err
	}

	// Initialize connection to database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Match from language and set source option if provided
	options := translate.Options{Format: "text"}
	if !isEmptyString(from) {
		fromLang, err := matchLang(from, db)
		if err != nil {
			return nil, err
		}
		options.Source = fromLang
	}

	// Match to language to provided string
	toLang, err := matchLang(to, db)
	if err != nil {
		return nil, err
	}

	// Request translations
	translations, err := client.Translate(ctx, []string{text}, toLang, &options)
	if err != nil {
		return nil, err
	}

	// Create slice of translation strings to return to client
	stringTranslations := []string{}
	for _, translation := range translations {
		stringTranslations = append(stringTranslations, translation.Text)
	}
	return stringTranslations, nil
}

// Create the supported language database
func createLanguagesDB() error {
	// Remvoe existing database
	os.Remove(dbPath)

	// Initialize connection to database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create languages table if it does not exist
	createStmt := `create table if not exists languages (tag text not null, name text not null, date int not null);`
	if _, err := db.Exec(createStmt); err != nil {
		return fmt.Errorf("%q: %s", err, createStmt)
	}

	// Start database transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into languages(tag, name, date) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Initialize Cloud Translation client
	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Request supported languages from Cloud Translation API
	langs, err := client.SupportedLanguages(ctx, language.English)
	if err != nil {
		return err
	}

	// Insert each language into the database
	for _, lang := range langs {
		_, err = stmt.Exec(lang.Tag.String(), lang.Name, time.Now().Unix())
		if err != nil {
			return err
		}
	}

	// Commit changes
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Get supported languages from the database and return
// them as a slice of strings
func getSupportedLanguages() ([]string, error) {
	// Initialize connection to database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Verify database exists and is up to date
	if err := verifyDB(); err != nil {
		return nil, err
	}

	// Select the names column from the table
	rows, err := db.Query("select name from languages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Append each language to a slice and return it
	languages := []string{}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		languages = append(languages, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return languages, nil
}

// Check if a provided string is empty
func isEmptyString(s string) bool {
	return len(strings.Fields(s)) == 0
}

// Match a supported language string to a returned language tag
func matchLang(l string, db *sql.DB) (language.Tag, error) {
	// Formatted error with provided language
	invalidLang := fmt.Errorf("%w%s", ErrInvalidLang, l)

	// Query statement to check against the tag and name in the table
	rows, err := db.Query(fmt.Sprintf("select tag from languages where tag=\"%s\" or name=\"%s\" collate nocase", l, l))
	if err != nil {
		return language.Und, invalidLang
	}
	defer rows.Close()

	// Select the matching string tag from the table,
	// parse into a language tag and return it.
	var tag language.Tag
	for rows.Next() {
		var stringTag string
		err = rows.Scan(&stringTag)
		if err != nil {
			return language.Und, invalidLang
		}
		tag, err = language.Parse(stringTag)
		if err != nil {
			return language.Und, invalidLang
		}
	}

	if err := rows.Err(); err != nil {
		return language.Und, invalidLang
	}

	return tag, nil
}

// Check if the supported language database exists and create it
// if not. Also checks if the database is up to date, and if not,
// creates a new database.
func verifyDB() error {
	// Check if the database exists
	if _, err := os.Stat(dbPath); errors.Is(err, fs.ErrNotExist) {
		return createLanguagesDB()
	}

	// Initialize connection to database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Check if the database is outdated
	rows, err := db.Query("select date from languages limit 1")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var date int64
		err = rows.Scan(&date)
		if err != nil {
			return err
		}

		if date+(24*3600) < time.Now().Unix() {
			return createLanguagesDB()
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
