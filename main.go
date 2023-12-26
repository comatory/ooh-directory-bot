package main

import (
	"internal/bot"
	"internal/client"
	"internal/parser"
	"internal/processor"
	"internal/scraper"
	"log"
)

const URL = "https://ooh.directory/random/"

func main() {
	botConfig, botConfigError := bot.ReadConfiguration()

	if botConfigError != nil {
		log.Fatal(botConfigError)
	}

	httpClient := client.CreateHttpClient()
	html, err := scraper.ScrapeRandom(URL, &httpClient)

	if err != nil {
		log.Fatal(err)
	}

	results, parseError := parser.ParseResults(html)

	if parseError != nil {
		log.Println(parseError)
	}

	result, processError := processor.ProcessResultForAPI(&results, &processor.FileStorage{})

	if processError != nil {
		log.Println(processError)
	}

	postError := bot.PostResult(result, &botConfig, &httpClient)

	if postError != nil {
		log.Fatal(postError)
	}
}
