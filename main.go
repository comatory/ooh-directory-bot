package main

import (
	"fmt"
	"net/http"
)

const URL = "https://ooh.directory/random/"

func main() {
	httpClient := http.Client{}
	fmt.Println(ScrapeRandom(URL, httpClient))
}
