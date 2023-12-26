package bot

import (
	"errors"
	"os"
	"path"
	"strings"
)

const configFileName = ".env"

type Config struct {
	AccessToken  string
	BotServerUrl string
}

func isValidConfig(config *Config) bool {
	return config.AccessToken != "" && config.BotServerUrl != ""
}

func ReadConfiguration() (Config, error) {
	config := Config{}

	path := path.Join(".", configFileName)
	contents, err := os.ReadFile(path)

	if err != nil {
		return config, err
	}

	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		pair := strings.Split(line, "=")

		if pair[0] == "access_token" {
			config.AccessToken = pair[1]
		}

		if pair[0] == "bot_server_url" {
			config.BotServerUrl = pair[1]
		}
	}

	if !isValidConfig(&config) {
		return config, errors.New("Invalid config")
	}

	return config, nil
}
