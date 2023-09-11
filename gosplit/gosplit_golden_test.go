package gosplit_test

import (
	"inaz2/GoSplit/gosplit"
	"flag"
	"testing"

	"github.com/tenntenn/golden"
)

var (
	flagUpdate bool
)

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

func testByLinesGoldenTarget(dir string) error {
	filePath := "testdata/example.txt"
	prefix := "TestByLines-"
	nLines := 10

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(dir)
	err := g.ByLines(nLines)
	if err != nil {
		return err
	}

	return nil
}

func testByNumberGoldenTarget(dir string) error {
	filePath := "testdata/example.txt"
	prefix := "TestByNumber-"
	nNumber := 4

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(dir)
	err := g.ByNumber(nNumber)
	if err != nil {
		return err
	}

	return nil
}

func testByBytesGoldenTarget(dir string) error {
	filePath := "testdata/example.txt"
	prefix := "TestByBytes-"
	nBytes := int64(512)

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(dir)
	err := g.ByBytes(nBytes)
	if err != nil {
		return err
	}

	return nil
}

func TestByLinesGolden(t *testing.T) {
	dir := t.TempDir()
	if err := testByLinesGoldenTarget(dir); err != nil {
		t.Fatal("unexpected error:", err)
	}

	got := golden.Txtar(t, dir)
	if diff := golden.Check(t, flagUpdate, "testdata", "TestByLinesGolden", got); diff != "" {
		t.Error(diff)
	}
}

func TestByNumberGolden(t *testing.T) {
	dir := t.TempDir()
	if err := testByNumberGoldenTarget(dir); err != nil {
		t.Fatal("unexpected error:", err)
	}

	got := golden.Txtar(t, dir)
	if diff := golden.Check(t, flagUpdate, "testdata", "TestByNumberGolden", got); diff != "" {
		t.Error(diff)
	}
}

func TestByBytesGolden(t *testing.T) {
	dir := t.TempDir()
	if err := testByBytesGoldenTarget(dir); err != nil {
		t.Fatal("unexpected error:", err)
	}

	got := golden.Txtar(t, dir)
	if diff := golden.Check(t, flagUpdate, "testdata", "TestByBytesGolden", got); diff != "" {
		t.Error(diff)
	}
}
