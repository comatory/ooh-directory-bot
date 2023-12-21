package main

type Result struct {
	url        string
	title      string
	summary    string
	authorName string
}

func ParseResults(body string) []Result {
	return make([]Result, 0)
}
