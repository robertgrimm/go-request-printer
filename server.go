package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	port = 8080
)

type RequestPrinter struct{}

// Print out the request.  It assumes it's a JSON payload.  If it
// fails to parse it as JSON, it just prints it as is.
func (rp *RequestPrinter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handling %s %s request from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
	fmt.Println("Headers:")
	for k, v := range r.Header {
		fmt.Printf("\t%s:\t%s\n", k, v)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading body")
		return
	}

	fmt.Println("Body:")
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "  ")
	if err == nil {
		fmt.Println(string(prettyJSON.Bytes()))
	} else {
		fmt.Println(string(body))
	}

	fmt.Println()
}

func main() {
	handler := &RequestPrinter{}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
