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

type languageFile struct {
	Date      int
	Languages []translate.Language
}

var supportedLanguages []translate.Language
var languagesPath = "languages.json"
var ErrInvalidLang = errors.New("invalid or unsupported language: ")
var ErrInvalidRequest = errors.New("invalid to language or text to translate")

func TranslateText(from, to, text string) (string, error) {
	if len(strings.Fields(to)) == 0 || len(strings.Fields(text)) == 0 {
		return "", ErrInvalidRequest
	}

	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	supportedLanguages, err = checkLanguages(&ctx, client)
	if err != nil {
		return "", err
	}

	options := translate.Options{Format: "text"}
	if len(strings.Fields(from)) != 0 {
		fromName, err := matchNameToLang(from, supportedLanguages)
		if err != nil {
			fromTag, err := matchTagToLang(from, supportedLanguages)
			if err != nil {
				return "", err
			}
			options.Source = fromTag
		} else {
			options.Source = fromName
		}
	}

	var translateTo language.Tag
	toName, err := matchNameToLang(to, supportedLanguages)
	if err != nil {
		toTag, err := matchTagToLang(to, supportedLanguages)
		if err != nil {
			return "", err
		}
		translateTo = toTag
	} else {
		translateTo = toName
	}

	translation, err := client.Translate(ctx, []string{text}, translateTo, &options)
	if err != nil {
		return "", err
	}

	return translation[0].Text, nil
}

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

func matchTagToLang(l string, langs []translate.Language) (language.Tag, error) {
	for _, lang := range langs {
		if strings.EqualFold(lang.Tag.String(), l) {
			return lang.Tag, nil
		}
	}
	return language.Und, fmt.Errorf("%w%s", ErrInvalidLang, l)
}

func matchNameToLang(l string, langs []translate.Language) (language.Tag, error) {
	for _, lang := range langs {
		if strings.EqualFold(lang.Name, l) {
			return lang.Tag, nil
		}
	}
	return language.Und, fmt.Errorf("%w%s", ErrInvalidLang, l)
}
