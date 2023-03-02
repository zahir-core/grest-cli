package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	downloadStoplightElementsJS()
	downloadStoplightElementsCSS()
}

func downloadStoplightElementsJS() {
	res, err := http.Get("https://unpkg.com/@stoplight/elements/web-components.min.js")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	filename := "openapi-ui/stoplight-elements-web-components.min.js"
	os.Remove(filename)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func downloadStoplightElementsCSS() {
	res, err := http.Get("https://unpkg.com/@stoplight/elements/styles.min.css")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	filename := "openapi-ui/stoplight-elements-styles.min.css"
	os.Remove(filename)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		log.Fatal(err)
	}
}
