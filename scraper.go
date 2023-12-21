package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func RequestRandomsHtmlBody(url string, client http.Client) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "ooh-directory-random-bot")
	req.Header.Set("Accept-Language", "en-us, en-gb, en")
	req.Header.Set("Accept", "text/html")

	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Request failed: %d", res.StatusCode))
	}

	bytes, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ScrapeRandom(url string, client http.Client) (string, error) {
	return RequestRandomsHtmlBody(url, client)
}
