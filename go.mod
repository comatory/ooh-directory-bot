module comatory/ooh-directory-bot

go 1.21.5

require internal/parser v0.0.0

replace internal/parser => ./internal/parser

require internal/scraper v0.0.0

replace internal/scraper => ./internal/scraper

require internal/processor v0.0.0

require golang.org/x/net v0.19.0 // indirect

replace internal/processor => ./internal/processor
