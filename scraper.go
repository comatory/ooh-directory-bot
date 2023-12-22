package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func prepareRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "ooh-directory-random-bot")
	req.Header.Set("Accept-Language", "en-us, en-gb, en")
	req.Header.Set("Accept", "text/html")

	return req, nil
}

func performRequest(req *http.Request, client *http.Client) (*http.Response, error) {
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Request failed: %d", res.StatusCode))
	}

	return res, nil
}

func readResponseBody(res *http.Response) (string, error) {
	bytes, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ScrapeRandom(url string, client http.Client) (string, error) {
	req, prepareErr := prepareRequest(url)

	if prepareErr != nil {
		return "", prepareErr
	}

	res, requestError := performRequest(req, &client)

	if requestError != nil {
		return "", requestError
	}

	body, readError := readResponseBody(res)

	if readError != nil {
		return "", readError
	}

	return body, nil
}
