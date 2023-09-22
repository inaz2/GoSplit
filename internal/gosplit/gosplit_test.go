package gosplit_test

import (
	"inaz2/GoSplit/internal/gosplit"

	"bufio"
	"bytes"
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
		{prefix + "aa", 10},
		{prefix + "ab", 10},
		{prefix + "ac", 10},
		{prefix + "ad", 10},
		{prefix + "ae", 2},
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

func TestByLines_EmptyFile(t *testing.T) {
	t.Parallel()

	filePath := "testdata/empty"
	prefix := "TestByLines_EmptyFile-"
	outDir := t.TempDir()
	nLines := 10

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	gerr := g.ByLines(nLines)
	if gerr != nil {
		t.Fatal("ByLines() failed:", gerr)
	}

	outFileName := prefix + "aa"
	outFilePath := path.Join(outDir, outFileName)
	_, err := os.Stat(outFilePath)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", outFilePath)
	}
}

func TestByLines_Directory(t *testing.T) {
	t.Parallel()

	filePath := "testdata"
	prefix := "TestByLines_Directory-"
	outDir := t.TempDir()
	nLines := 10

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByLines(nLines)
	if err == nil {
		t.Errorf("ByLines() with a directory should be error")
	}
}

func TestByLines_InvalidNLines(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestByLines_InvalidNLines-"
	outDir := t.TempDir()
	nLines := 0

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByLines(nLines)
	if err == nil {
		t.Errorf("ByLines(%#v) should be error", nLines)
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
		{prefix + "aa", 363},
		{prefix + "ab", 363},
		{prefix + "ac", 363},
		{prefix + "ad", 366},
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

func TestByNumber_EmptyFile(t *testing.T) {
	t.Parallel()

	filePath := "testdata/empty"
	prefix := "TestByNumber_EmptyFile-"
	outDir := t.TempDir()
	nNumber := 4
	outFiles := []struct {
		name   string
		nBytes int64
	}{
		{prefix + "aa", 0},
		{prefix + "ab", 0},
		{prefix + "ac", 0},
		{prefix + "ad", 0},
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

func TestByNumber_Directory(t *testing.T) {
	t.Parallel()

	filePath := "testdata"
	prefix := "TestByNumber_Directory-"
	outDir := t.TempDir()
	nNumber := 4

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByNumber(nNumber)
	if err == nil {
		t.Errorf("ByNumber() with a directory should be error")
	}
}

func TestByNumber_Stdin(t *testing.T) {
	t.Parallel()

	filePath := "-"
	prefix := "TestByNumber_Stdin-"
	outDir := t.TempDir()
	nNumber := 4

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByNumber(nNumber)
	if err == nil {
		t.Errorf("ByNumber() with STDIN should be error")
	}
}

func TestByNumber_InvalidNBytes(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestByNumber_InvalidNBytes-"
	outDir := t.TempDir()
	nNumber := 0

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByNumber(nNumber)
	if err == nil {
		t.Errorf("ByNumber(%#v) should be error", nNumber)
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
		{prefix + "aa", 512},
		{prefix + "ab", 512},
		{prefix + "ac", 431},
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

func TestByBytes_EmptyFile(t *testing.T) {
	t.Parallel()

	filePath := "testdata/empty"
	prefix := "TestByBytes_EmptyFile-"
	outDir := t.TempDir()
	nBytes := int64(512)

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	gerr := g.ByBytes(nBytes)
	if gerr != nil {
		t.Fatal("ByBytes() failed:", gerr)
	}

	outFileName := prefix + "aa"
	outFilePath := path.Join(outDir, outFileName)
	_, err := os.Stat(outFilePath)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", outFilePath)
	}
}

func TestByBytes_Directory(t *testing.T) {
	t.Parallel()

	filePath := "testdata"
	prefix := "TestByBytes_Directory-"
	outDir := t.TempDir()
	nBytes := int64(512)

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByBytes(nBytes)
	if err == nil {
		t.Errorf("ByBytes() with a directory should be error")
	}
}

func TestByBytes_InvalidNBytes(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestByBytes_InvalidNBytes-"
	outDir := t.TempDir()
	nBytes := int64(0)

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	err := g.ByBytes(nBytes)
	if err == nil {
		t.Errorf("ByBytes(%#v) should be error", nBytes)
	}
}

func TestSetVerboseWriter(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestSetVerboseWriter-"
	outDir := t.TempDir()
	nNumber := 4

	path := path.Join(outDir, prefix)
	want := "creating file \"" + path + "aa\"\n" +
		"creating file \"" + path + "ab\"\n" +
		"creating file \"" + path + "ac\"\n" +
		"creating file \"" + path + "ad\"\n"

	var b bytes.Buffer
	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	g.SetVerboseWriter(&b)
	err := g.ByNumber(nNumber)
	if err != nil {
		t.Fatal("ByNumber() failed:", err)
	}

	got := b.String()
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestSetNumericSuffix(t *testing.T) {
	t.Parallel()

	filePath := "testdata/example.txt"
	prefix := "TestSetNumericSuffix-"
	outDir := t.TempDir()
	nNumber := 4
	outFiles := []struct {
		name   string
		nBytes int64
	}{
		{prefix + "00", 363},
		{prefix + "01", 363},
		{prefix + "02", 363},
		{prefix + "03", 366},
	}

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	g.SetNumericSuffix(true)
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

func TestSetElideEmptyFiles(t *testing.T) {
	t.Parallel()

	filePath := "testdata/empty"
	prefix := "TestSetElideEmptyFiles-"
	outDir := t.TempDir()
	nNumber := 4

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(outDir)
	g.SetElideEmptyFiles(true)
	gerr := g.ByNumber(nNumber)
	if gerr != nil {
		t.Fatal("ByNumber() failed:", gerr)
	}

	outFileName := prefix + "aa"
	outFilePath := path.Join(outDir, outFileName)
	_, err := os.Stat(outFilePath)
	if err == nil {
		t.Errorf("os.Stat(%#v) should be error", outFilePath)
	}
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
