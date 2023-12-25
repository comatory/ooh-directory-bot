package processor

import (
	"os"
	"path"
	"bufio"
	"internal/parser"
)

const RecordFileName = "records.txt"

type Storage interface {
	getRecord() (*os.File, error)
	storeRecord(result *parser.Result, file *os.File) error
	filterOutPreviousResults(results *[]parser.Result, file *os.File) []parser.Result
}

type FileStorage struct {}

func (*FileStorage) getRecord() (*os.File, error) {
	path := path.Join(".", RecordFileName)

	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func (*FileStorage) storeRecord(result *parser.Result, file *os.File) error {
	_, err := file.WriteString(result.Url + "\n")

	return err
}

func (*FileStorage) filterOutPreviousResults(results *[]parser.Result, file *os.File) []parser.Result {
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
