package gosplit_test

import (
	"inaz2/GoSplit/gosplit"

	"errors"
	"fmt"
	"strings"
	"testing"
)

func goSplitError1() error {
	return gosplit.GoSplitErrorf("root error")
}

func goSplitError2() error {
	if err := goSplitError1(); err != nil {
		return gosplit.GoSplitErrorf("nested error: %w", err)
	}
	return nil
}

func TestGoSplitError(t *testing.T) {
	t.Parallel()

	want := struct {
		message   string
		rootFrame string
	}{
		message:   "root error",
		rootFrame: "gosplit_test.goSplitError1",
	}

	err := goSplitError1()
	message := fmt.Sprintf("%v", err)
	detailed := fmt.Sprintf("%+v", err)
	if message != want.message {
		t.Errorf("message = %#v, want %#v", message, want.message)
	}
	if ok := strings.HasPrefix(detailed, want.message); !ok {
		t.Errorf("detailed = %#v, want HasPrefix(detailed, %#v)", detailed, want.message)
	}
	if strings.Count(detailed, want.rootFrame) != 1 {
		t.Errorf("detailed = %#v, want Count(detailed, %#v) == 1", detailed, want.rootFrame)
	}
}

func TestGoSplitErrorNested(t *testing.T) {
	t.Parallel()

	want := struct {
		message   string
		rootFrame string
	}{
		message:   "nested error: root error",
		rootFrame: "gosplit_test.goSplitError1",
	}

	err := goSplitError2()
	message := fmt.Sprintf("%v", err)
	detailed := fmt.Sprintf("%+v", err)
	if message != want.message {
		t.Errorf("message = %#v, want %#v", message, want.message)
	}
	if ok := strings.HasPrefix(detailed, want.message); !ok {
		t.Errorf("detailed = %#v, want HasPrefix(detailed, %#v)", detailed, want.message)
	}
	if strings.Count(detailed, want.rootFrame) != 1 {
		t.Errorf("detailed = %#v, want Count(detailed, %#v) == 1", detailed, want.rootFrame)
	}
}

func TestGoSplitErrorIs(t *testing.T) {
	t.Parallel()

	want := errors.New("an error")

	err := gosplit.GoSplitErrorf("TestGoSplitErrorIs: %w", want)
	if ok := errors.Is(err, gosplit.GoSplitErr); !ok {
		t.Errorf("errors.Is() should return true")
	}
}

func TestGoSplitErrorUnwrap(t *testing.T) {
	t.Parallel()

	want := errors.New("an error")

	err := gosplit.GoSplitErrorf("TestGoSplitErrorIs: %w", want)
	if ok := errors.Is(err, want); !ok {
		t.Errorf("errors.Is() should return true")
	}
}
