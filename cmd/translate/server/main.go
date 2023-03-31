package main

import (
	"log"
	"net/http"
)

func main() {
	createDB()
	startCron()
	http.HandleFunc("/json", http.HandlerFunc(jsonHandler))
	http.HandleFunc("/languages", http.HandlerFunc(languagesHandler))
	http.HandleFunc("/translate", http.HandlerFunc(queryHandler))
	log.Println(http.ListenAndServe(":8080", nil))
}
