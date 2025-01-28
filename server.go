package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

const (
	port        = 8080
	tlsPort     = 8443
	tlsCertFile = "server.crt"
	tlsKeyFile  = "server.key"
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading body")
		return
	}

	fmt.Println("=== Start of body ===")
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "  ")
	if err == nil {
		fmt.Println(string(prettyJSON.Bytes()))
	} else {
		fmt.Println(string(body))
	}
	fmt.Println("=== End of body ===")
	fmt.Println()
}

func getCurrentDir() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("Couldn't find go-request-printer source directory")
	}
	return filepath.Dir(currentFile)
}

func main() {
	handler := &RequestPrinter{}

	fmt.Printf("Listening on port %d (http) and %d (https)\n", port, tlsPort)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		curDir := getCurrentDir()
		certFile := path.Join(curDir, tlsCertFile)
		keyFile := path.Join(curDir, tlsKeyFile)
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			fmt.Printf("error loading certificate %s, %s: %v", certFile, keyFile, err)
			os.Exit(1)
		}
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
		tlsServer := &http.Server{
			Addr:      fmt.Sprintf(":%d", tlsPort),
			Handler:   handler,
			TLSConfig: tlsConfig,
		}
		log.Fatal(tlsServer.ListenAndServeTLS("", ""))
	}()

	wg.Wait()
}
