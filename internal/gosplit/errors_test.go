package gosplit_test

import (
	"inaz2/GoSplit/internal/gosplit"

	"errors"
	"fmt"
	"strings"
	"testing"
)

var errRoot = errors.New("root error")

func errGoSplit1() error {
	return gosplit.GoSplitErrorf("failed to something: %w", errRoot)
}

func errGoSplit2() error {
	if err := errGoSplit1(); err != nil {
		return gosplit.GoSplitErrorf("failed to errGoSplit1: %w", err)
	}
	return nil
}

func TestErrGoSplit(t *testing.T) {
	t.Parallel()

	err := errGoSplit1()
	cases := map[string]struct {
		in           string
		want         string
		expectPrefix bool
	}{
		"%v":   {"%v", "failed to something: root error", false},
		"%+v":  {"%+v", "failed to something: root error\n", true},
		"%#v":  {"%#v", "&gosplit.GeneralizedError{kind: &errors.errorString{s:\"gosplit\"}, err: &fmt.wrapError{msg:\"failed to something: root error\", ", true},
		"%#+v": {"%#+v", "&gosplit.GeneralizedError{kind: &errors.errorString{s:\"gosplit\"}, err: &fmt.wrapError{msg:\"failed to something: root error\", ", true},
		"%s":   {"%s", "failed to something: root error", false},
		"%q":   {"%q", "\"failed to something: root error\"", false},
		"%x":   {"%x", "6661696c656420746f20736f6d657468696e673a20726f6f74206572726f72", false},
		"%X":   {"%X", "6661696C656420746F20736F6D657468696E673A20726F6F74206572726F72", false},
		"%d":   {"%d", "&{%!d(string=failed to something: root error) ", true},
		"%Z":   {"%Z", "&{%!Z(string=failed to something: root error) ", true},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := fmt.Sprintf(tt.in, err)
			if tt.expectPrefix {
				if ok := strings.HasPrefix(got, tt.want); !ok {
					t.Errorf("HasPrefix(%#v, %#v) = false, want true", got, tt.want)
				}
			} else {
				if tt.want != got {
					t.Errorf("fmt.Sprintf(%#v) = %#v, want %#v", tt.in, got, tt.want)
				}
			}
		})
	}
}

func TestErrGoSplitStack(t *testing.T) {
	t.Parallel()

	err := errGoSplit2()
	frames := []string{"gosplit_test.errGoSplit2", "gosplit_test.errGoSplit1"}

	detailed := fmt.Sprintf("%+v", err)
	goReprDetailed := fmt.Sprintf("%#+v", err)

	for _, frame := range frames {
		if got := strings.Count(detailed, frame); got != 1 {
			t.Errorf("Count(%#v, %#v) = %#v, want 1", detailed, frame, got)
		}
		if got := strings.Count(goReprDetailed, frame); got != 1 {
			t.Errorf("Count(%#v, %#v) = %#v, want 1", goReprDetailed, frame, got)
		}
	}
}

func TestErrGoSplitIs(t *testing.T) {
	t.Parallel()

	err := errGoSplit1()

	if ok := errors.Is(errRoot, gosplit.ErrGoSplit); ok {
		t.Errorf("errors.Is(errRoot, gosplit.ErrGoSplit) = true, want false")
	}
	if ok := errors.Is(err, gosplit.ErrGoSplit); !ok {
		t.Errorf("errors.Is(err, gosplit.ErrGoSplit) = false, want true")
	}
	if ok := errors.Is(err, errRoot); !ok {
		t.Errorf("errors.Is(err, errRoot) = false, want true")
	}
}
