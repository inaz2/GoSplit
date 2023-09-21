package gerrors_test

import (
	"inaz2/GoSplit/internal/gerrors"

	"errors"
	"fmt"
	"strings"
	"testing"
)

var errGErrors = errors.New("gerrors")

func GErrorsErrorf(format string, a ...any) error {
	return gerrors.Errorf(errGErrors, format, a...)
}

var errRoot = errors.New("root error")

func errorTest1() error {
	return GErrorsErrorf("failed to something: %w", errRoot)
}

func errorTest2() error {
	if err := errorTest1(); err != nil {
		return GErrorsErrorf("failed to errorTest1: %w", err)
	}
	return nil
}

func TestErrGoSplit(t *testing.T) {
	t.Parallel()

	err := errorTest1()
	cases := map[string]struct {
		in           string
		want         string
		expectPrefix bool
	}{
		"%v":   {"%v", "failed to something: root error", false},
		"%+v":  {"%+v", "failed to something: root error\n", true},
		"%#v":  {"%#v", "&gerrors.generalizedError{kind: &errors.errorString{s:\"gerrors\"}, err: &fmt.wrapError{msg:\"failed to something: root error\", ", true},
		"%#+v": {"%#+v", "&gerrors.generalizedError{kind: &errors.errorString{s:\"gerrors\"}, err: &fmt.wrapError{msg:\"failed to something: root error\", ", true},
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

	err := errorTest2()
	frames := []string{"errors_test.errorTest2", "errors_test.errorTest1"}

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

	err := errorTest1()

	if ok := errors.Is(errRoot, errGErrors); ok {
		t.Errorf("errors.Is(errRoot, errGErrors) = true, want false")
	}
	if ok := errors.Is(err, errGErrors); !ok {
		t.Errorf("errors.Is(err, errGErrors) = false, want true")
	}
	if ok := errors.Is(err, errRoot); !ok {
		t.Errorf("errors.Is(err, errRoot) = false, want true")
	}
}
