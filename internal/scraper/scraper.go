package scraper

import (
	"internal/client"
	"io"
	"net/http"
)

func readResponseBody(res *http.Response) (string, error) {
	bytes, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ScrapeRandom(url string, client client.HttpClient) (string, error) {
	req, prepareErr := client.NewRequest(url, http.MethodGet)

	if prepareErr != nil {
		return "", prepareErr
	}

	res, requestError := client.DispatchRequest(req)

	if requestError != nil {
		return "", requestError
	}

	body, readError := readResponseBody(res)

	if readError != nil {
		return "", readError
	}

	return body, nil
}
