package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

const serviceURL = "https://translate-jmmhwrz5fq-uc.a.run.app/json"

func main() {
	from := flag.String("from", "", "language to translate from")
	to := flag.String("to", "", "language to translate to")
	flag.Parse()

	if *to == "" {

	}

	var text string
	for _, arg := range flag.Args() {
		text += arg
	}

	if text == "" {
		flag.Usage()
		return
	}

	bs, err := json.Marshal(map[string]string{"from": *from, "to": *to, "text": text})
	if err != nil {

	}
	req, err := http.NewRequest(http.MethodPut, serviceURL, bytes.NewReader(bs))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

	}
	defer resp.Body.Close()

	rs, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("translate \"%s\" to %s: %s\n", text, *to, string(rs))
}
