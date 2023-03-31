package translate

import (
	"context"
	"errors"
	"strings"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// Custom errors that can be returned joined with other errors
var (
	// no translation language provided
	ErrNoTo = errors.New("no translation language provided")
	// no translation text provided
	ErrNoText = errors.New("no translation text provided")
	// error initializing GCP Translation client:
	ErrInitClient = errors.New("error initializing GCP Translation client: ")
	// error parsing from language tag:
	ErrParseFrom = errors.New("error parsing from language tag: ")
	// error parsing to language tag:
	ErrParseTo = errors.New("error parsing to language tag: ")
	// error translating with GCP Translation API:
	ErrTranslate = errors.New("error translating with GCP Translation API: ")
	// error getting supported languages from GCP Translation API:
	ErrGetLangs = errors.New("error getting supported languages from GCP Translation API: ")
)

// Translation contains all pertinent information for a translation
type Translation struct {
	from       string
	to         string
	original   string
	translated string
}

// Translate original text to another language
// If from is not provided, the source language will be assumed
// Returns a Translation struct and an error
func TranslateText(from, to, original string) (Translation, error) {

	// Translation to return
	var translation Translation

	// Check if to or text are empty before sending a translation request
	if isEmptyString(to) {
		return translation, ErrNoTo
	}
	if isEmptyString(original) {
		return translation, ErrNoText
	}
	translation.original = original

	// Initialize Cloud Translation client
	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		return translation, errors.Join(ErrInitClient, err)
	}
	defer client.Close()

	// If provided, match from language string to language tag
	// and setting source option
	options := translate.Options{Format: "text"}
	if !isEmptyString(from) {
		fromLang, err := language.Parse(from)
		if err != nil {
			return translation, errors.Join(ErrParseFrom, err)
		}
		options.Source = fromLang
		translation.from = from
	}

	// Match provided to language string to language tag
	toLang, err := language.Parse(to)
	if err != nil {
		return translation, errors.Join(ErrParseTo, err)
	}
	translation.to = to

	// Request translations
	translations, err := client.Translate(ctx, []string{original}, toLang, &options)
	if err != nil {
		return translation, errors.Join(ErrTranslate, err)
	}

	// If from not provided, assign detected language source
	if isEmptyString(translation.from) {
		translation.from = translations[0].Source.String()
	}
	translation.translated = translations[0].Text
	return translation, nil
}

// Get supported languages from Cloud Translation API.
// Returns a map of supported language names:tags as strings.
func GetSupportedLanguages() (map[string]string, error) {
	// Initialize Cloud Translation client
	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		return nil, errors.Join(ErrInitClient, err)
	}
	defer client.Close()

	// Request supported languages from Cloud Translation API
	langs, err := client.SupportedLanguages(ctx, language.English)
	if err != nil {
		return nil, errors.Join(ErrGetLangs, err)
	}

	// Return supported languages as name:tag pairs
	// Ex: english:en
	languages := map[string]string{}
	for _, lang := range langs {
		languages[lang.Name] = lang.Tag.String()
	}

	return languages, nil
}

// Check if a provided string is empty
func isEmptyString(s string) bool {
	return len(strings.Fields(s)) == 0
}
