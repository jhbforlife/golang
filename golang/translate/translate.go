package translate

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

var supportedLanguages []translate.Language
var errInvalidLang = errors.New("%s language not supported or is invalid")

func TranslateText(from, to, text string) (string, error) {
	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	supportedLanguages, err = getSupportedLanguages(&ctx, client)
	if err != nil {
		return "", err
	}

	options := translate.Options{}
	if len(strings.Fields(from)) != 0 {
		fromName, err := matchNameToLang(from, supportedLanguages)
		if err != nil {
			fromTag, err := matchTagToLang(from, supportedLanguages)
			if err != nil {
				return "", fmt.Errorf(err.Error(), from)
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
			return "", fmt.Errorf(err.Error(), to)
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

func getSupportedLanguages(ctx *context.Context, client *translate.Client) ([]translate.Language, error) {
	langs, err := client.SupportedLanguages(*ctx, language.English)
	if err != nil {
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
	return language.Und, errInvalidLang
}

func matchNameToLang(l string, langs []translate.Language) (language.Tag, error) {
	for _, lang := range langs {
		if strings.EqualFold(lang.Name, l) {
			return lang.Tag, nil
		}
	}
	return language.Und, errInvalidLang
}
