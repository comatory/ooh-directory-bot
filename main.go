package main

import (
	"log"
	"net/http"
	"internal/parser"
	"internal/scraper"
	"internal/processor"
)

const URL = "https://ooh.directory/random/"

func main() {
	httpClient := http.Client{}
	html, err := scraper.ScrapeRandom(URL, httpClient)

	if err != nil {
		log.Fatal(err)
	}

	results, parseError := parser.ParseResults(html)

	if parseError != nil {
		log.Println(parseError)
	}

	result, processError := processor.ProcessResultForAPI(&results)

	if processError != nil {
		log.Println(processError)
	}

	log.Print(result)
}
