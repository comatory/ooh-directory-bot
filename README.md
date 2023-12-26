# ooh-directory-bot

Mastodon bot which posts a single random post from [ooh.directory](https://ooh.directory).

## Usage

Create `.env` configuration file next to binary. You can also create this file by copying `.env.example` to `.env` file from the source.

This configuration file must contain `access_token` (obtained from your server's application dashboard) and `bot_server_url`, which is the URL of the server, e.g. `https://botsin.space`.

Running the binary will then scrape the site and post to given bot's server URL `statuses` endpoint. The application keeps track of posted URLs in text file `records.txt` which will be written in binary's location.

The frequency of posting depends on how often you run the binary, use scheduler such as `cron` to set it up.

## Build

Run `make build` to obtain binaries for Windows/Linux/macOS. They will be in `outputs/` directory.

## Test

Run `make test` to run the test suites.

## Development

Run `make install` to copy pre-commit hook for formatting.
