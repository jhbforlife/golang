package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestJSONHandler(t *testing.T) {
	tc := []map[string]string{
		{"from": "english", "to": "french", "text": "hello", "expected": "200"},
		{"from": "", "to": "fr", "text": "hello", "expected": "200"},
		{"from": "eng", "to": "french", "text": "hello", "expected": "400"},
	}
	for _, c := range tc {
		t.Run(fmt.Sprintf("f:%s.t:%s.e:%s", c["from"], c["to"], c["expected"]), func(t *testing.T) {
			b, err := json.Marshal(c)
			if err != nil {
				t.Errorf("expected nil json error, got: %v", err)
			}
			req := httptest.NewRequest(http.MethodPut, "/json", bytes.NewBuffer(b))
			req.Header.Set("content-type", "application/json")
			w := httptest.NewRecorder()
			jsonHandler(w, req)
			if c["expected"] != strconv.Itoa(w.Result().StatusCode) {
				t.Errorf("got %v, want %v", w.Result().StatusCode, c["expected"])
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	tc := []*http.Request{
		httptest.NewRequest(http.MethodGet, "/translate?from=english&to=french&text=hello", nil),
	}
	w := httptest.NewRecorder()
	getHandler(w, tc[0])
	if w.Result().StatusCode != 200 {
		t.Errorf("got %v, want 200", w.Result().StatusCode)
	}
}
