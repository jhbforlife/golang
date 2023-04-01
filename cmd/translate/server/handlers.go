package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/jhbforlife/golang/translate"
)

// Incoming JSON request format
type translateRequest struct {
	Source, ToLang, Original string
}

// Handle incoming JSON POST requests
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.ErrNotSupported.ErrorString, http.StatusMethodNotAllowed)
		return
	}

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

	translateAndWrite(w, &req)
}

// Handle incoming GET requests for supported languages
func languagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.ErrNotSupported.ErrorString, http.StatusMethodNotAllowed)
		return
	}

	langs, err := getSupportedLanguages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonLangs, err := json.Marshal(langs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(jsonLangs)
}

// Handle incoming GET requests with query parameters
func queryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.ErrNotSupported.ErrorString, http.StatusMethodNotAllowed)
		return
	}
	vars := r.URL.Query()
	if !vars.Has("to") || !vars.Has("original") {
		http.Error(w, "missing 'to' or 'original'", http.StatusBadRequest)
		return
	}

	var req translateRequest
	if vars.Has("source") {
		req.Source = vars.Get("source")
	}
	req.ToLang = vars.Get("to")
	req.Original = vars.Get("original")

	translateAndWrite(w, &req)
}

// Requests translation using translate package and writes response to client
func translateAndWrite(w http.ResponseWriter, req *translateRequest) {
	source, err := matchLang(req.Source)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Source = source

	to, err := matchLang(req.ToLang)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.ToLang = to

	translation, err := matchTranslation(*req)
	if err != nil {
		newTranslation, err := translate.TranslateText(req.Source, req.ToLang, req.Original)
		if err != nil {
			if errors.Is(err, translate.ErrNoTo) || errors.Is(err, translate.ErrNoText) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}
		translation = newTranslation
	}

	jsonTranslation, err := json.Marshal(translation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(jsonTranslation)

	if err := insertTranslationIntoTable(translation); err != nil {
		log.Println(err)
	}
}
