package main

import (
	"io"
	"net/http"
)

type Result struct {
	url        string
	title      string
	summary    string
	authorName string
}

// func parseResults(body string) []Result {
// }

func RequestRandomsHtmlBody(url string, client http.Client) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", "ooh-directory-random-bot")
	req.Header.Set("Accept-Language", "en-us, en-gb, en")
	req.Header.Set("Accept", "text/html")

	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}

	bytes, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func ScrapeRandom(url string, client http.Client) string {
	htmlBody := RequestRandomsHtmlBody(url, client)

	return htmlBody
}
