package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elazarl/goproxy"
)

type URL struct {
	Path        string          `json:"path"`
	Query       string          `json:"query"`
	Response        json.RawMessage `json:"response"`
}

type URLS struct{
	Host  string  `json:"host"`
	Urls []URL `json:"urls"`
}

func main() {
	fmt.Println("Listening on port 8080")
	if len(os.Args) < 2 {
		log.Fatal("Please provide the path to the JSON file")
	}

	jsonFilePath := os.Args[1]

	data, err := readURLsFromFile(jsonFilePath)
	if err != nil {
		log.Fatal("Error reading URLs:", err)
	}

	proxy := goproxy.NewProxyHttpServer()

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if ctx.Req.URL.Path == "/stop" {
			os.Exit(0)
		}

        if ctx.Req.URL.Host == data.Host {
        fmt.Println(" ")
        fmt.Println("Intercepted:", ctx.Req.URL.Host)

		    for _, url := range data.Urls {
			    if  strings.Contains(ctx.Req.URL.Path, url.Path) && ctx.Req.URL.RawQuery == url.Query {
				    fmt.Println("Modified:", url.Path, url.Query)
					fmt.Println(" ")

				    jsonResponse := url.Response
				    resp.Body = io.NopCloser(bytes.NewReader(jsonResponse))
				    resp.ContentLength = int64(len(jsonResponse))
				    resp.Header.Set("Content-Type", "application/json")
			    }
			}
		}

		return resp
	})

	log.Fatal(http.ListenAndServe(":8080", proxy))
}

func readURLsFromFile(filename string) (URLS, error) {
	file, err := os.Open(filename)
	if err != nil {
		return URLS{}, err
	}
	defer file.Close()

	var urls URLS
	err = json.NewDecoder(file).Decode(&urls)
	if err != nil {
		return URLS{}, err
	}

	return urls, nil
}
