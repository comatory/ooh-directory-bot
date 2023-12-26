package processor

import (
	"internal/parser"
)

type ResultNotFoundError struct{}

func (*ResultNotFoundError) Error() string {
	return "Result was not selected"
}

func ProcessResultForAPI(results *[]parser.Result, storage Storage) (*parser.Result, error) {
	file, fileError := storage.GetRecord()

	defer file.Close()

	if fileError != nil {
		empty := parser.NewEmptyResult()
		return &empty, fileError
	}

	scanner := storage.ReadRecord(file)

	filteredResults := storage.FilterOutPreviousResults(results, scanner)

	if len(filteredResults) == 0 {
		empty := parser.NewEmptyResult()
		return &empty, &ResultNotFoundError{}
	}

	firstResult := filteredResults[0]

	// TODO: This should be moved when bot succesfully posts to API
	storeError := storage.StoreRecord(&firstResult, file)

	if storeError != nil {
		return &firstResult, storeError
	}

	return &firstResult, nil
}
