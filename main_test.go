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

func TestSplitByNumber(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestSplitByNumber-"
	nNumber := 4
	outFiles := []struct {
		name   string
		nBytes int64
	}{
		{"TestSplitByNumber-aa", 363},
		{"TestSplitByNumber-ab", 363},
		{"TestSplitByNumber-ac", 363},
		{"TestSplitByNumber-ad", 366},
	}

	err := splitByNumber(filePath, prefix, nNumber)
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

func TestSplitByNumberEmpty(t *testing.T) {
	filePath := "testdata/empty"
	prefix := "TestSplitByNumberEmpty-"
	nNumber := 4

	err := splitByNumber(filePath, prefix, nNumber)
	if err != nil {
		t.Fatal(err)
	}

	fileName := "TestSplitByNumberEmpty-aa"
	_, err = os.Stat(fileName)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", fileName)
	}
}

func TestSplitByNumberStdin(t *testing.T) {
	filePath := "-"
	prefix := "TestSplitByNumberEmpty-"
	nNumber := 4

	err := splitByNumber(filePath, prefix, nNumber)
	if err == nil {
		t.Errorf("splitByNumber() with stdin should be error")
	}
}

func TestSplitByNumberEmptyPrefix(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := ""
	nNumber := 4

	err := splitByNumber(filePath, prefix, nNumber)
	if err == nil {
		t.Errorf("empty prefix should be error")
	}
}

func TestSplitByNumberInvalidNBytes(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestSplitByNumberInvalidNBytes-"
	nNumber := 0

	err := splitByNumber(filePath, prefix, nNumber)
	if err == nil {
		t.Errorf("non-positive nNumber should be error")
	}
}

func TestSplitByBytes(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestSplitByBytes-"
	nBytes := int64(512)
	outFiles := []struct {
		name   string
		nBytes int64
	}{
		{"TestSplitByBytes-aa", 512},
		{"TestSplitByBytes-ab", 512},
		{"TestSplitByBytes-ac", 431},
	}

	err := splitByBytes(filePath, prefix, nBytes)
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

func TestSplitByBytesEmpty(t *testing.T) {
	filePath := "testdata/empty"
	prefix := "TestSplitByBytesEmpty-"
	nBytes := int64(512)

	err := splitByBytes(filePath, prefix, nBytes)
	if err != nil {
		t.Fatal(err)
	}

	fileName := "TestSplitByBytesEmpty-aa"
	_, err = os.Stat(fileName)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", fileName)
	}
}

func TestSplitByBytesEmptyPrefix(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := ""
	nBytes := int64(512)

	err := splitByBytes(filePath, prefix, nBytes)
	if err == nil {
		t.Errorf("empty prefix should be error")
	}
}

func TestSplitByBytesInvalidNBytes(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestSplitByBytesInvalidNBytes-"
	nBytes := int64(0)

	err := splitByBytes(filePath, prefix, nBytes)
	if err == nil {
		t.Errorf("non-positive nBytes should be error")
	}
}
