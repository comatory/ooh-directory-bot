package main

import (
	"bufio"
	"os"
	"path"
)

const recordFileName = "records.txt"

type ResultNotFoundError struct{}

func (*ResultNotFoundError) Error() string {
	return "Result was not selected"
}

func getRecordFile() (*os.File, error) {
	path := path.Join(".", recordFileName)

	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func filterOutPreviousResults(results *[]Result, file *os.File) []Result {
	scanner := bufio.NewScanner(file)
	list := make([]Result, len(*results))

	copy(list, *results)

	for scanner.Scan() {
		for index, result := range *results {
			if result.url == scanner.Text() {
				list = append(list[:index], list[index+1:]...)
			}
		}
	}

	return list
}

func storeRecord(result *Result, file *os.File) error {
	_, err := file.WriteString(result.url + "\n")

	return err
}

func ProcessResultForAPI(results *[]Result) (*Result, error) {
	file, fileError := getRecordFile()

	defer file.Close()

	if fileError != nil {
		empty := NewEmptyResult()
		return &empty, fileError
	}

	filteredResults := filterOutPreviousResults(results, file)

	if len(filteredResults) == 0 {
		empty := NewEmptyResult()
		return &empty, &ResultNotFoundError{}
	}

	firstResult := filteredResults[0]

	storeError := storeRecord(&firstResult, file)

	if storeError != nil {
		return &firstResult, storeError
	}

	return &firstResult, nil
}
