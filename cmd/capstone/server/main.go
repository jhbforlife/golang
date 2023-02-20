package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/jhbatshipt/golang/translate"
)

type translateRequest struct {
	From, To, Text string
}

func main() {
	http.HandleFunc("/json", http.HandlerFunc(jsonHandler))
	http.HandleFunc("/translate", http.HandlerFunc(getHandler))
	panic(http.ListenAndServe(":8080", nil))
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var req translateRequest
	if err := json.Unmarshal(resp, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	translations, err := translate.TranslateText(req.From, req.To, req.Text)
	text := translations[0]
	if err != nil {
		if errors.Is(err, translate.ErrInvalidLang) || errors.Is(err, translate.ErrNoText) || errors.Is(err, translate.ErrNoToLang) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}

	if _, err := w.Write([]byte(text)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	if !vars.Has("to") || !vars.Has("text") {
		http.Error(w, "missing 'to' or 'text'", http.StatusBadRequest)
		return
	}

	var req translateRequest
	if vars.Has("from") {
		req.From = vars.Get("from")
	}
	req.To = vars.Get("to")
	req.Text = vars.Get("text")

	translations, err := translate.TranslateText(req.From, req.To, req.Text)
	text := translations[0]
	if err != nil {
		if errors.Is(err, translate.ErrInvalidLang) || errors.Is(err, translate.ErrNoText) || errors.Is(err, translate.ErrNoToLang) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}
	if _, err := w.Write([]byte(text)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
