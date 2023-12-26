package bot

import (
	"bytes"
	"encoding/json"
	"internal/client"
	"internal/parser"
	"net/http"
)

const endpoint = "/api/v1/statuses"

type Payload struct {
	Status string `json:"status"`
}

func addAuthorization(req *http.Request, config *Config) {
	req.Header.Set("Authorization", "Bearer "+config.AccessToken)
}

func addJsonContentType(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

func createPayload(result *parser.Result) (*bytes.Buffer, error) {
	status := result.Url + " " + result.Title

	if result.HasAuthorName() {
		status += " (by " + result.AuthorName

		if result.HasUpdatedAt() {
			status += ", " + result.FormatUpdatedAt() + ")"
		} else {
			status += ")"
		}
	}

	jsonValue, err := json.Marshal(Payload{
		Status: status,
	})

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonValue), nil
}

func PostResult(result *parser.Result, config *Config, client *client.Client) error {
	payload, payloadErr := createPayload(result)

	if payloadErr != nil {
		return payloadErr
	}

	req, prepareErr := client.NewRequestWithBody(config.BotServerUrl+endpoint, http.MethodPost, payload)

	addAuthorization(req, config)
	addJsonContentType(req)

	if prepareErr != nil {
		return prepareErr
	}

	_, requestError := client.DispatchRequest(req)

	if requestError != nil {
		return requestError
	}

	return nil
}
