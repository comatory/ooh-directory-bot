module comatory/ooh-directory-bot

go 1.23.2

require internal/parser v0.0.0

replace internal/parser => ./internal/parser

require internal/scraper v0.0.0

replace internal/scraper => ./internal/scraper

require (
	internal/bot v0.0.0-00010101000000-000000000000
	internal/client v0.0.0
	internal/storage v0.0.0
)

replace internal/client => ./internal/client

require golang.org/x/net v0.31.0 // indirect

replace internal/storage => ./internal/storage

replace internal/bot => ./internal/bot
