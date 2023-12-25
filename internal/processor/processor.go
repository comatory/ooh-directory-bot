package processor

import (
	"internal/parser"
)

type ResultNotFoundError struct{}

func (*ResultNotFoundError) Error() string {
	return "Result was not selected"
}


func ProcessResultForAPI(results *[]parser.Result, storage Storage) (*parser.Result, error) {
	file, fileError := storage.getRecord()

	defer file.Close()

	if fileError != nil {
		empty := parser.NewEmptyResult()
		return &empty, fileError
	}

	filteredResults := storage.filterOutPreviousResults(results, file)

	if len(filteredResults) == 0 {
		empty := parser.NewEmptyResult()
		return &empty, &ResultNotFoundError{}
	}

	firstResult := filteredResults[0]

	// TODO: This should be moved when bot succesfully posts to API
	storeError := storage.storeRecord(&firstResult, file)

	if storeError != nil {
		return &firstResult, storeError
	}

	return &firstResult, nil
}
