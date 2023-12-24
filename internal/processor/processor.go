package processor

import (
	"bufio"
	"internal/parser"
	"os"
	"path"
)

const RecordFileName = "records.txt"

type ResultNotFoundError struct{}

func (*ResultNotFoundError) Error() string {
	return "Result was not selected"
}

func getRecordFile() (*os.File, error) {
	path := path.Join(".", RecordFileName)

	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func filterOutPreviousResults(results *[]parser.Result, file *os.File) []parser.Result {
	scanner := bufio.NewScanner(file)
	list := make([]parser.Result, len(*results))

	copy(list, *results)

	for scanner.Scan() {
		for index, result := range *results {
			if result.Url == scanner.Text() {
				list = append(list[:index], list[index+1:]...)
			}
		}
	}

	return list
}

func storeRecord(result *parser.Result, file *os.File) error {
	_, err := file.WriteString(result.Url + "\n")

	return err
}

func ProcessResultForAPI(results *[]parser.Result) (*parser.Result, error) {
	file, fileError := getRecordFile()

	defer file.Close()

	if fileError != nil {
		empty := parser.NewEmptyResult()
		return &empty, fileError
	}

	filteredResults := filterOutPreviousResults(results, file)

	if len(filteredResults) == 0 {
		empty := parser.NewEmptyResult()
		return &empty, &ResultNotFoundError{}
	}

	firstResult := filteredResults[0]

	storeError := storeRecord(&firstResult, file)

	if storeError != nil {
		return &firstResult, storeError
	}

	return &firstResult, nil
}
