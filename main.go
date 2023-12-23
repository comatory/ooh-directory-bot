package main

import (
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

	results, parseError := ParseResults(html)

	if parseError != nil {
		log.Println(parseError)
	}

	result, processError := ProcessResultForAPI(&results)

	if processError != nil {
		log.Println(processError)
	}

	log.Print(result)
}
