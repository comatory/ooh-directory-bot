package main

import (
	"flag"
	"fmt"
	"internal/bot"
	"internal/client"
	"internal/parser"
	"internal/scraper"
	"internal/storage"
	"log"
	"os"
)

const URL = "https://ooh.directory/random/"

type Flags struct {
	configFilePath  *string
	recordsFilePath *string
}

func main() {
	log.SetOutput(os.Stdout)

	log.Println("Starting bot")

	flags := &Flags{
		configFilePath:  flag.String("config-file", "", "Path to the configuration file"),
		recordsFilePath: flag.String("records-file", "", "Path to the records file"),
	}

  flag.Parse()

	botConfig, botConfigError := bot.ReadConfiguration(flags.configFilePath)

	if botConfigError != nil {
		log.Fatal(botConfigError)
	}

	log.Println("Configuration loaded OK")

	httpClient := client.CreateHttpClient()
	fileStorage := storage.FileStorage{}

	loadError := fileStorage.Load(flags.recordsFilePath)

	defer fileStorage.Close()

	if loadError != nil {
		log.Fatal(loadError)
	}

	html, err := scraper.ScrapeRandom(URL, &httpClient)

	if err != nil {
		log.Fatal(err)
	}

	results, parseError := parser.ParseResults(html)
	log.Println(fmt.Sprintf("Scraped %d results", len(results)))

	if parseError != nil {
		log.Println(parseError)
	}

	result, processError := scraper.ProcessResultForAPI(&results, &fileStorage)

	if processError != nil {
		log.Println(processError)
	}

	log.Println(fmt.Sprintf("Selected result: %s", result.Url))

	postError := bot.PostResult(result, &botConfig, &httpClient)

	if postError != nil {
		log.Fatal(postError)
	}

	storeError := fileStorage.StoreRecord(result)

	if storeError != nil {
		log.Panic(storeError)
	}

	log.Println("Bot finished")
}
