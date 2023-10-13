package gerrors_test

import (
	g "inaz2/GoSplit/internal/gerrors"

	"errors"
	"fmt"
	"io/fs"
	"strings"
	"testing"
)

// usage in pkg/subpkg

// errSubPkg represents any errors in subpkg.
var errSubPkg = errors.New("subpkg")

// wrapperSubPkg is a error wrapper for subpkg.
var wrapperSubPkg = g.NewWrapper(errSubPkg)

// okSubPkg returns nil as a g.Error.
func okSubPkg() g.Error {
	return nil
}

// failSubPkg1 returns a g.Error linked to fs.ErrExist.
func failSubPkg1() g.Error {
	return wrapperSubPkg.Errorf("failed something in fs: %w", fs.ErrExist)
}

// failSubPkg2 returns a g.Error from failSubPkg1.
func failSubPkg2() g.Error {
	if err := failSubPkg1(); err != nil {
		return wrapperSubPkg.Errorf("failed to failSubPkg1: %w", err)
	}
	return nil
}

// usage in pkg

// errPkg represents any errors in pkg.
var errPkg = errors.New("pkg")

// errPkgInternal represents a specific error about subpkg.
var errPkgInternal = errors.New("failed something in subpkg")

// wrapperPkg is a error wrapper for pkg.
var wrapperPkg = g.NewWrapper(errPkg)

// failPkg returns a g.Error from failSubPkg2.
func failPkg() g.Error {
	if err := failSubPkg2(); err != nil {
		return wrapperPkg.Link(errPkgInternal, err)
	}
	return nil
}

func TestNil(t *testing.T) {
	t.Parallel()

	if err := okSubPkg(); err != nil {
		t.Errorf("got = %#v, want nil", err)
	}
}

func TestFormat(t *testing.T) {
	t.Parallel()

	err := failSubPkg1()
	cases := map[string]struct {
		in           string
		want         string
		expectPrefix bool
	}{
		"%v":   {"%v", "failed something in fs: file already exists", false},
		"%+v":  {"%+v", "failed something in fs: file already exists\n", true},
		"%#v":  {"%#v", "&gerrors.errorWithStack{err: &fmt.wrapErrors{msg:\"failed something in fs: file already exists\", ", true},
		"%#+v": {"%#+v", "&gerrors.errorWithStack{err: &fmt.wrapErrors{msg:\"failed something in fs: file already exists\", ", true},
		"%s":   {"%s", "failed something in fs: file already exists", false},
		"%q":   {"%q", "\"failed something in fs: file already exists\"", false},
		"%x":   {"%x", "6661696c656420736f6d657468696e6720696e2066733a2066696c6520616c726561647920657869737473", false},
		"%X":   {"%X", "6661696C656420736F6D657468696E6720696E2066733A2066696C6520616C726561647920657869737473", false},
		"%d":   {"%d", "&{%!d(string=failed something in fs: file already exists) ", true},
		"%Z":   {"%Z", "&{%!Z(string=failed something in fs: file already exists) ", true},
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

func TestFormat_Stack(t *testing.T) {
	t.Parallel()

	err := failSubPkg2()
	frames := []string{"errors_test.failSubPkg2", "errors_test.failSubPkg1"}

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

func TestIs(t *testing.T) {
	t.Parallel()

	err := failSubPkg1()
	if ok := errors.Is(err, errSubPkg); !ok {
		t.Errorf("errors.Is(err, errSubPkg) = false, want true")
	}
	if ok := errors.Is(err, fs.ErrExist); !ok {
		t.Errorf("errors.Is(err, fs.ErrExist) = false, want true")
	}
	if ok := errors.Is(errSubPkg, err); ok {
		t.Errorf("errors.Is(errSubPkg, err) = true, want false")
	}
}

func TestLink_Format(t *testing.T) {
	t.Parallel()

	err := failPkg()
	want := "failed something in subpkg"

	got := fmt.Sprintf("%v", err)
	if got != want {
		t.Errorf("got = %#v, want %#v", got, want)
	}
}

func TestLink_Format_Stack(t *testing.T) {
	t.Parallel()

	err := failPkg()
	frames := []string{"errors_test.failPkg", "errors_test.failSubPkg2", "errors_test.failSubPkg1"}

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

func TestLink_Is(t *testing.T) {
	t.Parallel()

	err := failPkg()

	if ok := errors.Is(err, errPkgInternal); !ok {
		t.Errorf("errors.Is(err, errPkgInternal) = false, want true")
	}
	if ok := errors.Is(err, errPkg); !ok {
		t.Errorf("errors.Is(err, errPkg) = false, want true")
	}
	if ok := errors.Is(err, errSubPkg); !ok {
		t.Errorf("errors.Is(err, errSubPkg) = false, want true")
	}
	if ok := errors.Is(err, fs.ErrExist); !ok {
		t.Errorf("errors.Is(err, fs.ErrExist) = false, want true")
	}
	if ok := errors.Is(errPkgInternal, err); ok {
		t.Errorf("errors.Is(errPkgInternal, err) = true, want false")
	}
}
