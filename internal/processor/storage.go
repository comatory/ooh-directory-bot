package processor

import (
	"bufio"
	"internal/parser"
	"os"
	"path"
)

const RecordFileName = "records.txt"

type Storage interface {
	GetRecord() (*os.File, error)
	ReadRecord(file *os.File) *bufio.Scanner
	StoreRecord(result *parser.Result, file *os.File) error
	FilterOutPreviousResults(results *[]parser.Result, scanner *bufio.Scanner) []parser.Result
}

type FileStorage struct{}

func (*FileStorage) GetRecord() (*os.File, error) {
	path := path.Join(".", RecordFileName)

	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func (*FileStorage) ReadRecord(file *os.File) *bufio.Scanner {
	return bufio.NewScanner(file)
}

func (*FileStorage) StoreRecord(result *parser.Result, file *os.File) error {
	_, err := file.WriteString(result.Url + "\n")

	return err
}

func (*FileStorage) FilterOutPreviousResults(results *[]parser.Result, scanner *bufio.Scanner) []parser.Result {
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
