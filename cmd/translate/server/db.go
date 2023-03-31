package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/jhbforlife/golang/translate"
)

// Path to supported languages database
const dbPath = "./languages.db"

// Custom errors that can be joined with other errors
var (
	// error opening database
	ErrOpenDB = errors.New("error opening database: ")
	// error creating table
	ErrCreateTable = errors.New("error creating table: ")
	// error starting transaction
	ErrStartTx = errors.New("error starting transaction: ")
	// error preparing transaction
	ErrPrepareTx = errors.New("error preparing transaction: ")
	// error executing statement
	ErrExecuteStmt = errors.New("error executing statement: ")
	// error committing transaction
	ErrCommitTx = errors.New("error committing transaction: ")
	// error querying table
	ErrQueryTable = errors.New("error querying table: ")
	// error scanning rows
	ErrScanRows = errors.New("error scanning rows: ")
	// error deleting table
	ErrDeleteTable = errors.New("error deleting table: ")
)

// Create the languages database
func createDB() error {
	// Initialize connection to database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Join(ErrOpenDB, err)
	}
	defer db.Close()

	// Create languages table if it does not exist
	langsStmt := `create table if not exists languages (name text not null, tag text not null);`
	if _, err := db.Exec(langsStmt); err != nil {
		return errors.Join(ErrCreateTable, err)
	}

	// Create translations table if it does not exist
	translationStmt := `create table if not exists translations (from text not null, to text not null, original text not null, translated text not null);`
	if _, err := db.Exec(translationStmt); err != nil {
		return errors.Join(ErrCreateTable, err)
	}

	// Get supported languages using translate package
	langs, err := translate.GetSupportedLanguages()
	if err != nil {
		return err
	}

	return insertLanguagesIntoTable(*db, langs)
}

// Get supported languages from the database and return
// them as a map[string]string with name:tag pairs
func getSupportedLanguages() (map[string]string, error) {
	// Verify database exists
	if err := verifyDB(); err != nil {
		return nil, err
	}

	// Initialize connection to database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.Join(ErrOpenDB, err)
	}
	defer db.Close()

	// Select the names column from the languages table
	rows, err := db.Query("select name, tag from languages")
	if err != nil {
		return nil, errors.Join(ErrQueryTable, err)
	}
	defer rows.Close()

	// Add each name:tag pair to a map and then return it
	languages := map[string]string{}
	for rows.Next() {
		var name string
		var tag string
		err = rows.Scan(&name, &tag)
		if err != nil {
			return nil, errors.Join(ErrScanRows, err)
		}
		languages["name"] = tag
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Join(ErrScanRows, err)
	}

	return languages, nil
}

// Insert language name:tag pairs into the languages table
func insertLanguagesIntoTable(db *sql.DB, langs map[string]string) error {
	// Verify database exists
	if err := verifyDB(); err != nil {
		log.Println(err)
	}

	// Start database transaction
	tx, err := db.Begin()
	if err != nil {
		return errors.Join(ErrStartTx, err)
	}
	stmt, err := tx.Prepare("insert into languages(name, tag) values(?, ?)")
	if err != nil {
		return errors.Join(ErrPrepareTx, err)
	}
	defer stmt.Close()

	// Insert each language into the database
	for name, tag := range langs {
		_, err = stmt.Exec(name, tag)
		if err != nil {
			return errors.Join(ErrExecuteStmt, err)
		}
	}

	// Commit changes
	err = tx.Commit()
	if err != nil {
		return errors.Join(ErrCommitTx, err)
	}

	return nil
}

// Insert a translation into the translations table
func insertTranslationIntoTable(t *translate.Translation) error {
	// Initialize connection to database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return errors.Join(ErrOpenDB, err)
	}
	defer db.Close()

	// Start database transaction
	tx, err := db.Begin()
	if err != nil {
		return errors.Join(ErrStartTx, err)
	}
	stmt, err := tx.Prepare("insert into translations(from, to, original, translated) values(?, ?, ?, ?)")
	if err != nil {
		return errors.Join(ErrPrepareTx, err)
	}
	defer stmt.Close()

	// Insert translation into the table
	_, err = stmt.Exec(t.from, t.to, t.original, t.translated)
	if err != nil {
		return errors.Join(ErrExecuteStmt, err)
	}

	// Commit changes
	err = tx.Commit()
	if err != nil {
		return errors.Join(ErrCommitTx, err)
	}

	return nil
}

// Match a provided language string against the name
// or tag in the languages table
func matchLang(l string) (string, error) {
	// Verify database exists
	if err := verifyDB(); err != nil {
		log.Println(err)
	}

	// Formatted error with provided language
	invalidLang := fmt.Errorf("invalid language: %s", l)

	// Initialize connection to database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", errors.Join(ErrOpenDB, err)
	}
	defer db.Close()

	// Query statement to check against the name and tag in the table
	rows, err := db.Query(fmt.Sprintf("select tag from languages where tag=\"%s\" or name=\"%s\" collate nocase", l, l))
	if err != nil {
		return "", invalidLang
	}
	defer rows.Close()

	// Select the matching tag from the table and return it
	var stringTag string
	for rows.Next() {
		err = rows.Scan(&stringTag)
		if err != nil {
			return "", invalidLang
		}
	}

	if err := rows.Err(); err != nil {
		return "", invalidLang
	}

	return stringTag, nil
}

// Delete all rows in the translations table using CRON
func resetTranslationsTable() {
	// Verify database exists
	if err := verifyDB(); err != nil {
		log.Println(err)
	}

	// Initialize connection to database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(errors.Join(ErrOpenDB, err))
	}
	defer db.Close()

	// Delete all rows in the translations table
	deleteStmt := "delete from translations"
	if _, err := db.Exec(deleteStmt); err != nil {
		log.Println(errors.Join(ErrDeleteTable, err))
	}
}

// Update the languages database using the translate package and CRON
func updateLanguagesTable() {
	// Verify database exists
	if err := verifyDB(); err != nil {
		log.Println(err)
	}

	// Initialize connection to database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(errors.Join(ErrOpenDB, err))
	}
	defer db.Close()

	// Get supported languages using translate package
	langs, err := translate.GetSupportedLanguages()
	if err != nil {
		log.Println(err)
	}

	// Delete all rows in the languages table
	deleteStmt := "delete from languages"
	if _, err := db.Exec(deleteStmt); err != nil {
		log.Println(errors.Join(ErrDeleteTable, err))
	}

	if err := insertLanguagesIntoTable(*db, langs); err != nil {
		log.Println(err)
	}
}

// Check if the database exists and create it if not.
func verifyDB() error {
	// Check if the database exists
	if _, err := os.Stat(dbPath); errors.Is(err, fs.ErrNotExist) {
		return createDB()
	}

	return nil
}
