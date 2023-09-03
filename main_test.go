package main

import (
	"bufio"
	"os"
	"testing"
)

func countLines(filePath string) (int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	nLines := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		nLines++
	}
	return nLines, nil
}

func countBytes(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	fileSize := fileInfo.Size()
	return fileSize, nil
}

func TestGoSplitByLines(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestGoSplitByLines-"
	nLines := 10
	outFiles := []struct {
		name   string
		nLines int
	}{
		{"TestGoSplitByLines-aa", 10},
		{"TestGoSplitByLines-ab", 10},
		{"TestGoSplitByLines-ac", 10},
		{"TestGoSplitByLines-ad", 10},
		{"TestGoSplitByLines-ae", 2},
	}

	gosplit := GoSplit{filePath, prefix}
	err := gosplit.ByLines(nLines)
	if err != nil {
		t.Fatal(err)
	}

	for _, outFile := range outFiles {
		result, err := countLines(outFile.name)
		if err != nil {
			t.Fatal(err)
		}
		if result != outFile.nLines {
			t.Errorf("countLines(%#v) = %#v, want %#v", outFile.name, result, outFile.nLines)
		}
		defer os.Remove(outFile.name)
	}
}

func TestGoSplitByLinesEmpty(t *testing.T) {
	filePath := "testdata/empty"
	prefix := "TestGoSplitByLinesEmpty-"
	nLines := 10

	gosplit := GoSplit{filePath, prefix}
	err := gosplit.ByLines(nLines)
	if err != nil {
		t.Fatal(err)
	}

	fileName := "TestGoSplitByLinesEmpty-aa"
	_, err = os.Stat(fileName)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", fileName)
	}
}

func TestGoSplitByLinesEmptyPrefix(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := ""
	nLines := 10

	gosplit := GoSplit{filePath, prefix}
	err := gosplit.ByLines(nLines)
	if err == nil {
		t.Errorf("empty prefix should be error")
	}
}

func TestGoSplitByLinesInvalidNLines(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestGoSplitByLinesInvalidNLines-"
	nLines := 0

	gosplit := GoSplit{filePath, prefix}
	err := gosplit.ByLines(nLines)
	if err == nil {
		t.Errorf("non-positive nLines should be error")
	}
}

func TestGoSplitByNumber(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestGoSplitByNumber-"
	nNumber := 4
	outFiles := []struct {
		name   string
		nBytes int64
	}{
		{"TestGoSplitByNumber-aa", 363},
		{"TestGoSplitByNumber-ab", 363},
		{"TestGoSplitByNumber-ac", 363},
		{"TestGoSplitByNumber-ad", 366},
	}

	gosplit := GoSplit{filePath, prefix}
	err := gosplit.ByNumber(nNumber)
	if err != nil {
		t.Fatal(err)
	}

	for _, outFile := range outFiles {
		result, err := countBytes(outFile.name)
		if err != nil {
			t.Fatal(err)
		}
		if result != outFile.nBytes {
			t.Errorf("countBytes(%#v) = %#v, want %#v", outFile.name, result, outFile.nBytes)
		}
		defer os.Remove(outFile.name)
	}
}

func TestGoSplitByNumberEmpty(t *testing.T) {
	filePath := "testdata/empty"
	prefix := "TestGoSplitByNumberEmpty-"
	nNumber := 4

	gosplit := GoSplit{filePath, prefix}
	err := gosplit.ByNumber(nNumber)
	if err != nil {
		t.Fatal(err)
	}

	fileName := "TestGoSplitByNumberEmpty-aa"
	_, err = os.Stat(fileName)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", fileName)
	}
}

func TestGoSplitByNumberStdin(t *testing.T) {
	filePath := "-"
	prefix := "TestGoSplitByNumberEmpty-"
	nNumber := 4

	gosplit := GoSplit{filePath, prefix}
	err := gosplit.ByNumber(nNumber)
	if err == nil {
		t.Errorf("ByNumber() with stdin should be error")
	}
}

func TestGoSplitByNumberEmptyPrefix(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := ""
	nNumber := 4

	gosplit := GoSplit{filePath, prefix}
	err := gosplit.ByNumber(nNumber)
	if err == nil {
		t.Errorf("empty prefix should be error")
	}
}

func TestGoSplitByNumberInvalidNBytes(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestGoSplitByNumberInvalidNBytes-"
	nNumber := 0

	gosplit := GoSplit{filePath, prefix}
	err := gosplit.ByNumber(nNumber)
	if err == nil {
		t.Errorf("non-positive nNumber should be error")
	}
}

func TestGoSplitByBytes(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestGoSplitByBytes-"
	nBytes := int64(512)
	outFiles := []struct {
		name   string
		nBytes int64
	}{
		{"TestGoSplitByBytes-aa", 512},
		{"TestGoSplitByBytes-ab", 512},
		{"TestGoSplitByBytes-ac", 431},
	}

	gosplit := GoSplit{filePath, prefix}
	err := gosplit.ByBytes(nBytes)
	if err != nil {
		t.Fatal(err)
	}

	for _, outFile := range outFiles {
		result, err := countBytes(outFile.name)
		if err != nil {
			t.Fatal(err)
		}
		if result != outFile.nBytes {
			t.Errorf("countBytes(%#v) = %#v, want %#v", outFile.name, result, outFile.nBytes)
		}
		defer os.Remove(outFile.name)
	}
}

func TestGoSplitByBytesEmpty(t *testing.T) {
	filePath := "testdata/empty"
	prefix := "TestGoSplitByBytesEmpty-"
	nBytes := int64(512)

	gosplit := GoSplit{filePath, prefix}
	err := gosplit.ByBytes(nBytes)
	if err != nil {
		t.Fatal(err)
	}

	fileName := "TestGoSplitByBytesEmpty-aa"
	_, err = os.Stat(fileName)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", fileName)
	}
}

func TestGoSplitByBytesEmptyPrefix(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := ""
	nBytes := int64(512)

	gosplit := GoSplit{filePath, prefix}
	err := gosplit.ByBytes(nBytes)
	if err == nil {
		t.Errorf("empty prefix should be error")
	}
}

func TestGoSplitByBytesInvalidNBytes(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestGoSplitByBytesInvalidNBytes-"
	nBytes := int64(0)

	gosplit := GoSplit{filePath, prefix}
	err := gosplit.ByBytes(nBytes)
	if err == nil {
		t.Errorf("non-positive nBytes should be error")
	}
}
