package main

import (
	"log"
	"fmt"
	"internal/bot"
	"internal/client"
	"internal/parser"
	"internal/processor"
	"internal/scraper"
)

const URL = "https://ooh.directory/random/"

func main() {
	log.Println("Starting bot")

	botConfig, botConfigError := bot.ReadConfiguration()

	if botConfigError != nil {
		log.Fatal(botConfigError)
	}

	log.Println("Configuration loaded OK")

	httpClient := client.CreateHttpClient()
	html, err := scraper.ScrapeRandom(URL, &httpClient)

	if err != nil {
		log.Fatal(err)
	}

	results, parseError := parser.ParseResults(html)
	log.Println(fmt.Sprintf("Scraped %d results", len(results)))

	if parseError != nil {
		log.Println(parseError)
	}

	result, processError := processor.ProcessResultForAPI(&results, &processor.FileStorage{})

	if processError != nil {
		log.Println(processError)
	}

	log.Println(fmt.Sprintf("Selected result: %s", result.Url))

	postError := bot.PostResult(result, &botConfig, &httpClient)

	if postError != nil {
		log.Fatal(postError)
	}

	log.Println("Bot finished")
}
