package storage

import (
	"bufio"
	"os"
	"path"
	"internal/parser"
)

const RecordFileName = "records.txt"

type Storage interface {
	Load() error
	GetFile() *os.File
	Close() error
	ReadRecord() *bufio.Scanner
	StoreRecord(result *parser.Result) error
	FilterOutPreviousResults(results *[]parser.Result, scanner *bufio.Scanner) []parser.Result
}

type FileStorage struct{
	file *os.File
}

func (storage *FileStorage) Load() error {
	path := path.Join(".", RecordFileName)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	storage.file = file

	return nil
}

func (storage *FileStorage) GetFile() *os.File {
	return storage.file
}

func (storage *FileStorage) Close() error {
	return storage.file.Close()
}

func (storage *FileStorage) ReadRecord() *bufio.Scanner {
	return bufio.NewScanner(storage.file)
}

func (storage *FileStorage) StoreRecord(result *parser.Result) error {
	_, err := storage.file.WriteString(result.Url + "\n")

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
