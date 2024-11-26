package scraper

import (
	"os"
	"io/fs"
	"testing"
	"testing/fstest"
	"bufio"
	"internal/parser"
	"internal/storage"
)

func getMockRecordsUrls(records *[]parser.Result) []byte {
	var urls []byte

	for _, record := range *records {
		urls = append(urls, []byte(record.Url + "\n")...)
	}

	return urls
}

type mockStorage struct {
	MockRecords []parser.Result
	file fs.File
}

func (storage *mockStorage) Load(recordFilePath *string) error {
	fs := fstest.MapFS{
		"test.txt": {
			Data: getMockRecordsUrls(&storage.MockRecords),
		},
	}

	file, err := fs.Open("test.txt")

	if err != nil {
		return err
	}

	storage.file = file

	return nil
}

func (*mockStorage) GetFile() *os.File {
	return nil
}

func (storage *mockStorage) ReadRecord() *bufio.Scanner {
	return bufio.NewScanner(storage.file)
}

func (*mockStorage) StoreRecord(result *parser.Result) error {
	return nil
}

func (*mockStorage) Close() error {
	return nil
}

func (*mockStorage) FilterOutPreviousResults(results *[]parser.Result, scanner *bufio.Scanner) []parser.Result {
	s := storage.FileStorage{}
	return s.FilterOutPreviousResults(results, scanner)
}

func TestProcessorWithFirstResult(t *testing.T) {
	results := []parser.Result{
		{
			Url: "https://ooh.directory/random/1",
		},
		{
			Url: "https://ooh.directory/random/2",
		},
	}

	storage := mockStorage{
		MockRecords: make([]parser.Result, 0),
	}
	storage.Load(nil)

	result, _ := ProcessResultForAPI(&results, &storage)

	if result.Url != "https://ooh.directory/random/1" {
	  t.Errorf("Expected %s, got %s", "https://ooh.directory/random/1", result.Url)
	}
}

func TestProcessorWithFilteredOutResult(t *testing.T) {
	results := []parser.Result{
		{
			Url: "https://ooh.directory/random/1",
		},
		{
			Url: "https://ooh.directory/random/2",
		},
	}

	mockRecords := []parser.Result{
		{
			Url: "https://ooh.directory/random/1",
		},
	}

	storage := mockStorage{
		MockRecords: mockRecords,
	}
	storage.Load(nil)

	result, _ := ProcessResultForAPI(&results, &storage)

	if result.Url != "https://ooh.directory/random/2" {
	  t.Errorf("Expected %s, got %s", "https://ooh.directory/random/2", result.Url)
	}
}

func TestProcessorWithoutAnyResult(t *testing.T) {
	results := []parser.Result{
		{
			Url: "https://ooh.directory/random/1",
		},
	}

	mockRecords := []parser.Result{
		{
			Url: "https://ooh.directory/random/1",
		},
	}

	storage := mockStorage{
		MockRecords: mockRecords,
	}
	storage.Load(nil)

	_, err := ProcessResultForAPI(&results, &storage)

	if err == nil {
	  t.Errorf("Expected error, got %v", err)
	}
}
