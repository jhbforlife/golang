package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const serviceURL = "" // no server currently live

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprint(w, "\nusage: translate [options] text\n\n-f, -from\n\tlanguage to translate from\n-h, -help\n\tprint this message\n-t, -to *required\n\tlanguage to translate to\n\n")
	}

	fFlag := flag.String("f", "", "language to translate from")
	fromFlag := flag.String("from", "", "language to translate from")
	hFlag := flag.Bool("h", false, "display help message")
	helpFlag := flag.Bool("help", false, "display help message")
	tFlag := flag.String("t", "", "language to translate to")
	toFlag := flag.String("to", "", "language to translate to")
	flag.Parse()

	if *fFlag == "" && *fromFlag == "" && !*hFlag && !*helpFlag && *tFlag == "" && *toFlag == "" {
		flag.Usage()
		return
	}

	if *hFlag || *helpFlag {
		flag.Usage()
		return
	}

	from := getStringFromFlags(*fFlag, *fromFlag)
	to := getStringFromFlags(*tFlag, *toFlag)

	if to == "" {
		fmt.Println("translate: no to language specified")
		flag.Usage()
		return
	}

	text := strings.Join(flag.Args(), " ")

	if len(strings.Fields(text)) == 0 {
		fmt.Println("translate: no text to translate")
		flag.Usage()
		return
	}

	bs, err := json.Marshal(map[string]string{"from": from, "to": to, "text": text})
	if checkErr(err) {
		return
	}

	req, err := http.NewRequest(http.MethodPut, serviceURL, bytes.NewReader(bs))
	if checkErr(err) {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if checkErr(err) {
		return
	}
	defer resp.Body.Close()

	rs, err := io.ReadAll(resp.Body)
	if checkErr(err) {
		return
	}

	fmt.Printf("translate \"%s\" to %s: %s\n", text, to, string(rs))
}

func getStringFromFlags(f, g string) string {
	if len(strings.Fields(f)) == 0 && len(strings.Fields(g)) != 0 {
		return g
	} else if len(strings.Fields(f)) != 0 && len(strings.Fields(g)) == 0 {
		return f
	} else {
		return ""
	}
}

func checkErr(err error) bool {
	if err != nil {
		fmt.Printf("translate: %v\n", err)
		return true
	}
	return false
}
