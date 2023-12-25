package processor

import (
	"testing"
	"testing/fstest"
	"bufio"
	"io"
	"internal/parser"
)

func TestWithNoPreviousResultsInStorage(t *testing.T) {
	fileMock := fstest.MapFS{
		"records.txt": &fstest.MapFile{
			Data: []byte(""),
		},
	}
	storage := FileStorage{}
	results := []parser.Result{
		{
			Url: "https://ooh.directory/random/1",
		},
		{
			Url: "https://ooh.directory/random/2",
		},
		{
			Url: "https://ooh.directory/random/3",
		},
	}

	mockFile, _ := fileMock.Open("records.txt")
	scanner := bufio.NewScanner(io.Reader(mockFile))
	filteredResults := storage.FilterOutPreviousResults(&results, scanner)

	if len(filteredResults) != 3 {
		t.Errorf("Expected \"%d\" results, got \"%d\"", 3, len(filteredResults))
	}
}

func TestWithPreviousResultsInStorage(t *testing.T) {
	fileMock := fstest.MapFS{
		"records.txt": &fstest.MapFile{
			Data: []byte("https://ooh.directory/random/2"),
		},
	}
	storage := FileStorage{}
	results := []parser.Result{
		{
			Url: "https://ooh.directory/random/1",
		},
		{
			Url: "https://ooh.directory/random/2",
		},
		{
			Url: "https://ooh.directory/random/3",
		},
	}

	mockFile, _ := fileMock.Open("records.txt")
	scanner := bufio.NewScanner(io.Reader(mockFile))
	filteredResults := storage.FilterOutPreviousResults(&results, scanner)

	if len(filteredResults) != 2 {
		t.Errorf("Expected \"%d\" results, got \"%d\"", 2, len(filteredResults))
	}

	if filteredResults[0].Url != "https://ooh.directory/random/1" {
		t.Errorf("Expected \"%s\" result, got \"%s\"", "https://ooh.directory/random/1", filteredResults[0].Url)
	}

	if filteredResults[1].Url != "https://ooh.directory/random/3" {
	  t.Errorf("Expected \"%s\" result, got \"%s\"", "https://ooh.directory/random/3", filteredResults[1].Url)
	}
}
