package gosplit_test

import (
	"inaz2/GoSplit/gosplit"

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

func TestParseSize(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestParseSize-"
	cases := map[string]struct {
		in        string
		want      int64
		expectErr bool
	}{
		"2":    {"2", 2, false},
		"2b":   {"2b", 2 * 512, false},
		"2K":   {"2K", 2 * 1024, false},
		"2KiB": {"2KiB", 2 * 1024, false},
		"2KB":  {"2KB", 2 * 1000, false},
		"7E":   {"7E", 7 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024, false},
		"9EB":  {"9EB", 9 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000, false},
		"8E":   {"8E", 0, true},
		"10EB": {"10EB", 0, true},
		"1Z":   {"1Z", 0, true},
		"1ZB":  {"1ZB", 0, true},
		"0":    {"0", 0, true},
		"0K":   {"0K", 0, true},
		"1.5":  {"1.5", 0, true},
		"-1":   {"-1", 0, true},
		"2iB":  {"2iB", 0, true},
		"2B":   {"2B", 0, true},
		"2biB": {"2biB", 0, true},
		"2bB":  {"2bB", 0, true},
		"X":    {"X", 0, true},
		"2X":   {"2X", 0, true},
		"2KX":  {"2KX", 0, true},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			g := gosplit.New(filePath, prefix)
			got, err := g.ParseSize(tt.in)
			if tt.expectErr && err == nil {
				t.Fatal("want err")
			}
			if !tt.expectErr && err != nil {
				t.Fatal("not want err:", err)
			}
			if tt.want != got {
				t.Errorf("ParseSize(%#v) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
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

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByLines(nLines)
	if err != nil {
		t.Fatal("ByLines() failed:", err)
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

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByLines(nLines)
	if err != nil {
		t.Fatal("ByLines() failed:", err)
	}

	fileName := "TestByLinesEmpty-aa"
	_, err = os.Stat(fileName)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", fileName)
	}
}

func TestByLinesInvalidNLines(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestByLinesInvalidNLines-"
	outDir := t.TempDir()
	nLines := 0

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByLines(nLines)
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

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByNumber(nNumber)
	if err != nil {
		t.Fatal("ByNumber() failed:", err)
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

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByNumber(nNumber)
	if err != nil {
		t.Fatal("ByNumber() failed:", err)
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

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByNumber(nNumber)
	if err == nil {
		t.Errorf("ByNumber() with stdin should be error")
	}
}

func TestByNumberInvalidNBytes(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestByNumberInvalidNBytes-"
	outDir := t.TempDir()
	nNumber := 0

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByNumber(nNumber)
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

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByBytes(nBytes)
	if err != nil {
		t.Fatal("ByBytes() failed:", err)
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

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByBytes(nBytes)
	if err != nil {
		t.Fatal("ByBytes() failed:", err)
	}

	fileName := "TestByBytesEmpty-aa"
	_, err = os.Stat(fileName)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", fileName)
	}
}

func TestByBytesInvalidNBytes(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestByBytesInvalidNBytes-"
	outDir := t.TempDir()
	nBytes := int64(0)

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByBytes(nBytes)
	if err == nil {
		t.Errorf("non-positive nBytes should be error")
	}
}
