package translate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// Data file for supported languages
// TODO: Change to database
type languageFile struct {
	Date      int
	Languages []translate.Language
}

// Slice of supported languages
var supportedLanguages []translate.Language

// TODO: Change to database
// Path to the stored languages file
var languagesPath = "languages.json"

// Error returned if a language is either invalid or not supported.
var ErrInvalidLang = errors.New("invalid or unsupported language: ")

// Error returned if no translation language is provided.
var ErrNoToLang = errors.New("no translation language provided")

// Error returned if no translation text is provided.
var ErrNoText = errors.New("no translation text provided")

// Returns a slice of supported language names.
func SupportedLanguages() ([]string, error) {
	// TODO: Implement
	return nil, nil
}

// Translate text to another language. Returns an error if either
// to or text are empty, or if either language is invalid or unsupported.
// If from is not provided, the source language will be assumed.
// Returns a slice of possible translations, ranked from hightest to lowest confidence.
func TranslateText(from, to, text string) ([]string, error) {

	// Check if to or text are empty before sending a translation request
	if isEmptyString(to) {
		return nil, ErrNoToLang
	}
	if isEmptyString(text) {
		return nil, ErrNoText
	}

	// Initialization
	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// TODO: Change to database
	// Retrieve supported languages
	supportedLanguages, err = checkLanguages(&ctx, client)
	if err != nil {
		return nil, err
	}

	// Match from language and set source option if provided
	options := translate.Options{Format: "text"}
	if !isEmptyString(from) {
		fromLang, err := matchLang(from, supportedLanguages)
		if err != nil {
			return nil, err
		}
		options.Source = fromLang
	}

	// Match to language to provided string
	toLang, err := matchLang(to, supportedLanguages)
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

// TODO: Change to database
func checkLanguages(ctx *context.Context, client *translate.Client) ([]translate.Language, error) {
	var languageFile languageFile
	bs, err := os.ReadFile(languagesPath)
	switch {
	case os.IsNotExist(err):
		langs, err := getSupportedLanguages(ctx, client)
		if err != nil {
			return nil, err
		}
		return langs, nil
	case err != nil:
		return nil, err
	}

	if err := json.Unmarshal(bs, &languageFile); err != nil {
		if err := os.Remove(languagesPath); err == nil {
			return checkLanguages(ctx, client)
		}
		return nil, err
	}
	if languageFile.Date != time.Now().Day() {
		langs, err := getSupportedLanguages(ctx, client)
		if err != nil {
			return nil, err
		}
		return langs, nil
	}
	return languageFile.Languages, nil
}

// TODO: Change to database
func getSupportedLanguages(ctx *context.Context, client *translate.Client) ([]translate.Language, error) {
	langs, err := client.SupportedLanguages(*ctx, language.English)
	if err != nil {
		return nil, err
	}
	langsToFile := languageFile{time.Now().Day(), langs}
	bs, err := json.Marshal(langsToFile)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(languagesPath, bs, 0444); err != nil {
		return nil, err
	}
	return langs, nil
}

// Check if a provided string is empty
func isEmptyString(s string) bool {
	return len(strings.Fields(s)) == 0
}

// Match a supported language string to a returned language tag
func matchLang(l string, langs []translate.Language) (language.Tag, error) {
	for _, lang := range langs {
		if strings.EqualFold(lang.Name, l) || strings.EqualFold(lang.Tag.String(), l) {
			return lang.Tag, nil
		}
	}
	return language.Und, fmt.Errorf("%w%s", ErrInvalidLang, l)
}
