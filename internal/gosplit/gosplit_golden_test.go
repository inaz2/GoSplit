package gosplit_test

import (
	"inaz2/GoSplit/internal/gosplit"

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

func targetTestByLines_Golden(dir string) error {
	filePath := "testdata/example.txt"
	prefix := "TestByLines_Golden-"
	nLines := 10

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(dir)
	err := g.ByLines(nLines)
	if err != nil {
		return err
	}

	return nil
}

func targetTestByNumber_Golden(dir string) error {
	filePath := "testdata/example.txt"
	prefix := "TestByNumber_Golden-"
	nNumber := 4

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(dir)
	err := g.ByNumber(nNumber)
	if err != nil {
		return err
	}

	return nil
}

func targetTestByNumber_EmptyFile_Golden(dir string) error {
	filePath := "testdata/empty"
	prefix := "TestByNumber_EmptyFile_Golden-"
	nNumber := 4

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(dir)
	err := g.ByNumber(nNumber)
	if err != nil {
		return err
	}

	return nil
}

func targetTestByBytes_Golden(dir string) error {
	filePath := "testdata/example.txt"
	prefix := "TestByBytes_Golden-"
	nBytes := int64(512)

	g := gosplit.New(filePath, prefix)
	g.SetOutDir(dir)
	err := g.ByBytes(nBytes)
	if err != nil {
		return err
	}

	return nil
}

func TestByLines_Golden(t *testing.T) {
	dir := t.TempDir()
	if err := targetTestByLines_Golden(dir); err != nil {
		t.Fatal("unexpected error:", err)
	}

	got := golden.Txtar(t, dir)
	if diff := golden.Check(t, flagUpdate, "testdata", "TestByLines_Golden", got); diff != "" {
		t.Error(diff)
	}
}

func TestByNumber_Golden(t *testing.T) {
	dir := t.TempDir()
	if err := targetTestByNumber_Golden(dir); err != nil {
		t.Fatal("unexpected error:", err)
	}

	got := golden.Txtar(t, dir)
	if diff := golden.Check(t, flagUpdate, "testdata", "TestByNumber_Golden", got); diff != "" {
		t.Error(diff)
	}
}

func TestByNumber_EmptyFile_Golden(t *testing.T) {
	dir := t.TempDir()
	if err := targetTestByNumber_EmptyFile_Golden(dir); err != nil {
		t.Fatal("unexpected error:", err)
	}

	got := golden.Txtar(t, dir)
	if diff := golden.Check(t, flagUpdate, "testdata", "TestByNumber_EmptyFile_Golden", got); diff != "" {
		t.Error(diff)
	}
}

func TestByBytes_Golden(t *testing.T) {
	dir := t.TempDir()
	if err := targetTestByBytes_Golden(dir); err != nil {
		t.Fatal("unexpected error:", err)
	}

	got := golden.Txtar(t, dir)
	if diff := golden.Check(t, flagUpdate, "testdata", "TestByBytes_Golden", got); diff != "" {
		t.Error(diff)
	}
}
