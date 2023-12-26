package processor

import (
	"internal/parser"
)

type ResultNotFoundError struct{}

func (*ResultNotFoundError) Error() string {
	return "Result was not selected"
}

func ProcessResultForAPI(results *[]parser.Result, storage Storage) (*parser.Result, error) {
	scanner := storage.ReadRecord()

	filteredResults := storage.FilterOutPreviousResults(results, scanner)

	if len(filteredResults) == 0 {
		empty := parser.NewEmptyResult()
		return &empty, &ResultNotFoundError{}
	}

	firstResult := filteredResults[0]

	return &firstResult, nil
}
