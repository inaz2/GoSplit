package gosplit_test

import (
	. "inaz2/GoSplit/gosplit"
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

func TestByLines(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestByLines-"
	nLines := 10
	outFiles := []struct {
		name   string
		nLines int
	}{
		{"TestByLines-aa", 10},
		{"TestByLines-ab", 10},
		{"TestByLines-ac", 10},
		{"TestByLines-ad", 10},
		{"TestByLines-ae", 2},
	}

	goSplit := NewGoSplit(filePath, prefix)
	err := goSplit.ByLines(nLines)
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

func TestByLinesEmpty(t *testing.T) {
	filePath := "testdata/empty"
	prefix := "TestByLinesEmpty-"
	nLines := 10

	goSplit := NewGoSplit(filePath, prefix)
	err := goSplit.ByLines(nLines)
	if err != nil {
		t.Fatal(err)
	}

	fileName := "TestByLinesEmpty-aa"
	_, err = os.Stat(fileName)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", fileName)
	}
}

func TestByLinesEmptyPrefix(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := ""
	nLines := 10

	goSplit := NewGoSplit(filePath, prefix)
	err := goSplit.ByLines(nLines)
	if err == nil {
		t.Errorf("empty prefix should be error")
	}
}

func TestByLinesInvalidNLines(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestByLinesInvalidNLines-"
	nLines := 0

	goSplit := NewGoSplit(filePath, prefix)
	err := goSplit.ByLines(nLines)
	if err == nil {
		t.Errorf("non-positive nLines should be error")
	}
}

func TestByNumber(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestByNumber-"
	nNumber := 4
	outFiles := []struct {
		name   string
		nBytes int64
	}{
		{"TestByNumber-aa", 363},
		{"TestByNumber-ab", 363},
		{"TestByNumber-ac", 363},
		{"TestByNumber-ad", 366},
	}

	goSplit := NewGoSplit(filePath, prefix)
	err := goSplit.ByNumber(nNumber)
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

func TestByNumberEmpty(t *testing.T) {
	filePath := "testdata/empty"
	prefix := "TestByNumberEmpty-"
	nNumber := 4

	goSplit := NewGoSplit(filePath, prefix)
	err := goSplit.ByNumber(nNumber)
	if err != nil {
		t.Fatal(err)
	}

	fileName := "TestByNumberEmpty-aa"
	_, err = os.Stat(fileName)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", fileName)
	}
}

func TestByNumberStdin(t *testing.T) {
	filePath := "-"
	prefix := "TestByNumberEmpty-"
	nNumber := 4

	goSplit := NewGoSplit(filePath, prefix)
	err := goSplit.ByNumber(nNumber)
	if err == nil {
		t.Errorf("ByNumber() with stdin should be error")
	}
}

func TestByNumberEmptyPrefix(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := ""
	nNumber := 4

	goSplit := NewGoSplit(filePath, prefix)
	err := goSplit.ByNumber(nNumber)
	if err == nil {
		t.Errorf("empty prefix should be error")
	}
}

func TestByNumberInvalidNBytes(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestByNumberInvalidNBytes-"
	nNumber := 0

	goSplit := NewGoSplit(filePath, prefix)
	err := goSplit.ByNumber(nNumber)
	if err == nil {
		t.Errorf("non-positive nNumber should be error")
	}
}

func TestByBytes(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestByBytes-"
	nBytes := int64(512)
	outFiles := []struct {
		name   string
		nBytes int64
	}{
		{"TestByBytes-aa", 512},
		{"TestByBytes-ab", 512},
		{"TestByBytes-ac", 431},
	}

	goSplit := NewGoSplit(filePath, prefix)
	err := goSplit.ByBytes(nBytes)
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

func TestByBytesEmpty(t *testing.T) {
	filePath := "testdata/empty"
	prefix := "TestByBytesEmpty-"
	nBytes := int64(512)

	goSplit := NewGoSplit(filePath, prefix)
	err := goSplit.ByBytes(nBytes)
	if err != nil {
		t.Fatal(err)
	}

	fileName := "TestByBytesEmpty-aa"
	_, err = os.Stat(fileName)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", fileName)
	}
}

func TestByBytesEmptyPrefix(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := ""
	nBytes := int64(512)

	goSplit := NewGoSplit(filePath, prefix)
	err := goSplit.ByBytes(nBytes)
	if err == nil {
		t.Errorf("empty prefix should be error")
	}
}

func TestByBytesInvalidNBytes(t *testing.T) {
	filePath := "testdata/example.txt"
	prefix := "TestByBytesInvalidNBytes-"
	nBytes := int64(0)

	goSplit := NewGoSplit(filePath, prefix)
	err := goSplit.ByBytes(nBytes)
	if err == nil {
		t.Errorf("non-positive nBytes should be error")
	}
}
