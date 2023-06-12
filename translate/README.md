# translate

`translate` is a simple package that leverages the [Go library](https://cloud.google.com/translate/docs/reference/libraries/v2/go) for the Google Cloud Translation API. See my [translate server](https://github.com/jhbforlife/golang/tree/main/cmd/translate/server) for an example implementation.

> #### _<sub>IMPORTANT</sub>_
>
> Authentication with an appropriate GCP project is required. This can be set up in a local environment using [gcloud](https://cloud.google.com/translate/docs/authentication#local-development) and is set up automatically when deploying to GCP. See [the docs](https://cloud.google.com/translate/docs/authentication) for more info.

## Available Methods

### _`TranslateText(source, to, original string) (Translation, error)`_
takes in three strings - the source language, the language to translate to, and the original text - and returns a [`Translation`](#available-structs) and an error. The `source` and `to` strings can be either the language name or the [language tag](https://www.w3.org/International/core/langtags/rfc3066bis.html), though the `source` language is not required.

### _`GetSupportedLanguages() (map[string]string, error)`_
returns the Cloud Translation API's supported languages as a map of name:tag string pairs and an error.

## Available Structs

```
type Translation struct {
	Source     string
	To         string
	Original   string
	Translated string
}
```