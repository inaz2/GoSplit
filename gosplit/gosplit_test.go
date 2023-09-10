package gosplit_test

import (
	. "inaz2/GoSplit/gosplit"
	"bufio"
	"os"
	"path"
	"testing"
)

func helperCountLines(t *testing.T, outDir string, filePath string) int {
	t.Helper()

	f, err := os.Open(path.Join(outDir, filePath))
	if err != nil {
		t.Fatal("failed to open:", err)
	}
	defer f.Close()

	nLines := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		nLines++
	}
	return nLines
}

func helperCountBytes(t *testing.T, outDir string, filePath string) int64 {
	t.Helper()

	fileInfo, err := os.Stat(path.Join(outDir, filePath))
	if err != nil {
		t.Fatal("failed to stat:", err)
	}

	fileSize := fileInfo.Size()
	return fileSize
}

func TestByLines(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestByLines-"
	outDir := t.TempDir()
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

	goSplit := NewGoSplit(filePath, prefix, outDir)
	err := goSplit.ByLines(nLines)
	if err != nil {
		t.Fatal("goSplit.ByLines() failed:", err)
	}

	for _, outFile := range outFiles {
		result := helperCountLines(t, outDir, outFile.name)
		if result != outFile.nLines {
			t.Errorf("helperCountLines(%#v) = %#v, want %#v", outFile.name, result, outFile.nLines)
		}
	}
}

func TestByLinesEmpty(t *testing.T) {
	t.Parallel()

	filePath := "testdata/empty"
	prefix := "TestByLinesEmpty-"
	outDir := t.TempDir()
	nLines := 10

	goSplit := NewGoSplit(filePath, prefix, outDir)
	err := goSplit.ByLines(nLines)
	if err != nil {
		t.Fatal("goSplit.ByLines() failed:", err)
	}

	fileName := "TestByLinesEmpty-aa"
	_, err = os.Stat(fileName)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", fileName)
	}
}

func TestByLinesEmptyPrefix(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := ""
	outDir := t.TempDir()
	nLines := 10

	goSplit := NewGoSplit(filePath, prefix, outDir)
	err := goSplit.ByLines(nLines)
	if err == nil {
		t.Errorf("empty prefix should be error")
	}
}

func TestByLinesInvalidNLines(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestByLinesInvalidNLines-"
	outDir := t.TempDir()
	nLines := 0

	goSplit := NewGoSplit(filePath, prefix, outDir)
	err := goSplit.ByLines(nLines)
	if err == nil {
		t.Errorf("non-positive nLines should be error")
	}
}

func TestByNumber(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestByNumber-"
	outDir := t.TempDir()
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

	goSplit := NewGoSplit(filePath, prefix, outDir)
	err := goSplit.ByNumber(nNumber)
	if err != nil {
		t.Fatal("goSplit.ByNumber() failed:", err)
	}

	for _, outFile := range outFiles {
		result := helperCountBytes(t, outDir, outFile.name)
		if result != outFile.nBytes {
			t.Errorf("helperCountBytes(%#v) = %#v, want %#v", outFile.name, result, outFile.nBytes)
		}
	}
}

func TestByNumberEmpty(t *testing.T) {
	t.Parallel()

	filePath := "testdata/empty"
	prefix := "TestByNumberEmpty-"
	outDir := t.TempDir()
	nNumber := 4

	goSplit := NewGoSplit(filePath, prefix, outDir)
	err := goSplit.ByNumber(nNumber)
	if err != nil {
		t.Fatal("goSplit.ByNumber() failed:", err)
	}

	fileName := "TestByNumberEmpty-aa"
	_, err = os.Stat(fileName)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", fileName)
	}
}

func TestByNumberStdin(t *testing.T) {
	t.Parallel()

	filePath := "-"
	prefix := "TestByNumberEmpty-"
	outDir := t.TempDir()
	nNumber := 4

	goSplit := NewGoSplit(filePath, prefix, outDir)
	err := goSplit.ByNumber(nNumber)
	if err == nil {
		t.Errorf("ByNumber() with stdin should be error")
	}
}

func TestByNumberEmptyPrefix(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := ""
	outDir := t.TempDir()
	nNumber := 4

	goSplit := NewGoSplit(filePath, prefix, outDir)
	err := goSplit.ByNumber(nNumber)
	if err == nil {
		t.Errorf("empty prefix should be error")
	}
}

func TestByNumberInvalidNBytes(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestByNumberInvalidNBytes-"
	outDir := t.TempDir()
	nNumber := 0

	goSplit := NewGoSplit(filePath, prefix, outDir)
	err := goSplit.ByNumber(nNumber)
	if err == nil {
		t.Errorf("non-positive nNumber should be error")
	}
}

func TestByBytes(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestByBytes-"
	outDir := t.TempDir()
	nBytes := int64(512)
	outFiles := []struct {
		name   string
		nBytes int64
	}{
		{"TestByBytes-aa", 512},
		{"TestByBytes-ab", 512},
		{"TestByBytes-ac", 431},
	}

	goSplit := NewGoSplit(filePath, prefix, outDir)
	err := goSplit.ByBytes(nBytes)
	if err != nil {
		t.Fatal("goSplit.ByBytes() failed:", err)
	}

	for _, outFile := range outFiles {
		result := helperCountBytes(t, outDir, outFile.name)
		if result != outFile.nBytes {
			t.Errorf("helperCountBytes(%#v) = %#v, want %#v", outFile.name, result, outFile.nBytes)
		}
	}
}

func TestByBytesEmpty(t *testing.T) {
	t.Parallel()

	filePath := "testdata/empty"
	prefix := "TestByBytesEmpty-"
	outDir := t.TempDir()
	nBytes := int64(512)

	goSplit := NewGoSplit(filePath, prefix, outDir)
	err := goSplit.ByBytes(nBytes)
	if err != nil {
		t.Fatal("goSplit.ByBytes() failed:", err)
	}

	fileName := "TestByBytesEmpty-aa"
	_, err = os.Stat(fileName)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", fileName)
	}
}

func TestByBytesEmptyPrefix(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := ""
	outDir := t.TempDir()
	nBytes := int64(512)

	goSplit := NewGoSplit(filePath, prefix, outDir)
	err := goSplit.ByBytes(nBytes)
	if err == nil {
		t.Errorf("empty prefix should be error")
	}
}

func TestByBytesInvalidNBytes(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestByBytesInvalidNBytes-"
	outDir := t.TempDir()
	nBytes := int64(0)

	goSplit := NewGoSplit(filePath, prefix, outDir)
	err := goSplit.ByBytes(nBytes)
	if err == nil {
		t.Errorf("non-positive nBytes should be error")
	}
}
