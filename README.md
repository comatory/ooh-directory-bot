# ooh-directory-bot

<img alt="logo for the bot depicting robot's head with surprised look" src="./assets/avatar.png" width="200" height="auto" />

Mastodon bot which posts a single random post from [ooh.directory](https://ooh.directory).
See the bot in action on [hachyderm.io instance](https://hachyderm.io/@ooh_directory_bot).

## Usage

### Pre-requisities

Create `.env` configuration file next to binary. You can also create this file by copying `.env.example` to `.env` file from the source. The location to this configuration file can be changed with flag `--config-file`.

This configuration file must contain `access_token` (obtained from your server's application dashboard) and `bot_server_url`, which is the URL of the server, e.g. `https://botsin.space`.

### Run

Running the binary will then scrape the site and post to given bot's server URL `statuses` endpoint. The application keeps track of posted URLs in text file `records.txt` which will be written in binary's location (otherwise configurable with `--records-file` flag).

The frequency of posting depends on how often you run the binary, use scheduler such as `cron` to set it up.

Logs are streamed to `stdout` so you can pipe them to a file if you want to.

## Build

Run `make build` to obtain binaries for Windows/Linux/macOS. They will be in `outputs/` directory.

## Test

Run `make test` to run the test suites.

## Development

Run `make install` to copy pre-commit hook for formatting.

## Deployment

It is recommended to use Docker. The included `Dockerfile` compiles the binary and runs it via cron, so need to set it up. **The interval is hard-coded** and will run the bot daily at 10 AM.

The docker setup assumes that you have folder `data/` in the same directory as the Dockerfile. This folder will contain `records.txt` file, you can also move existing one there. You also **must have** environment file placed in there and output log will be streamed to this folder as well.

Build the image first:

```bash
$ docker build -t ooh-directory-bot .
```

Then run the container (recommended settings for server-like environment):

```bash
$ docker run -d --restart always -v $(pwd)/data/:/app/data/ ooh-directory-bot
```

The downside of running the container this way with `cron` is that you won't be able to stop it with `CTRL+C` outside of daemon mode, you must stop the container with `docker stop <container_id>`.
