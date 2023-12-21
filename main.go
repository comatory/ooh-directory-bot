package main

import (
	"fmt"
	"log"
	"net/http"
)

const URL = "https://ooh.directory/random/"

func main() {
	httpClient := http.Client{}
	html, err := ScrapeRandom(URL, httpClient)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(html)
	fmt.Print(ParseResults(html))
}
