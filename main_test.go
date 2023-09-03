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

func TestSplitByLines(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestSplitByLines-"
	nLines := 10
	outFiles := []struct {
		name   string
		nLines int
	}{
		{"TestSplitByLines-aa", 10},
		{"TestSplitByLines-ab", 10},
		{"TestSplitByLines-ac", 10},
		{"TestSplitByLines-ad", 10},
		{"TestSplitByLines-ae", 2},
	}

	err := splitByLines(filePath, prefix, nLines)
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

func TestSplitByLinesEmpty(t *testing.T) {
	filePath := "testdata/empty"
	prefix := "TestSplitByLinesEmpty-"
	nLines := 10

	err := splitByLines(filePath, prefix, nLines)
	if err != nil {
		t.Fatal(err)
	}

	fileName := "TestSplitByLinesEmpty-aa"
	_, err = os.Stat(fileName)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", fileName)
	}
}

func TestSplitByLinesEmptyPrefix(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := ""
	nLines := 10

	err := splitByLines(filePath, prefix, nLines)
	if err == nil {
		t.Errorf("empty prefix should be error")
	}
}

func TestSplitByLinesInvalidNLines(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestSplitByLinesInvalidNLines-"
	nLines := 0

	err := splitByLines(filePath, prefix, nLines)
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
