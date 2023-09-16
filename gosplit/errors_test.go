package gosplit_test

import (
	"inaz2/GoSplit/gosplit"

	"errors"
	"fmt"
	"strings"
	"testing"
)

func errGoSplit1() error {
	return gosplit.GoSplitErrorf("root error")
}

func errGoSplit2() error {
	if err := errGoSplit1(); err != nil {
		return gosplit.GoSplitErrorf("nested error: %w", err)
	}
	return nil
}

func TestErrGoSplit(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in   string
		want string
	}{
		"%v":  {"%v", "root error"},
		"%#v": {"%#v", "&errors.errorString{s:\"root error\"}"},
		"%s":  {"%s", "root error"},
		"%q":  {"%q", "\"root error\""},
		"%x":  {"%x", "726f6f74206572726f72"},
		"%X":  {"%X", "726F6F74206572726F72"},
		"%d":  {"%d", "&{%!d(string=root error)}"},
		"%Z":  {"%Z", "&{%!Z(string=root error)}"},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := errGoSplit1()
			got := fmt.Sprintf(tt.in, err)
			if tt.want != got {
				t.Errorf("fmt.Sprintf(%#v) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestErrGoSplitStack(t *testing.T) {
	t.Parallel()

	want := struct {
		messagePrefix string
		goReprPrefix  string
		frames        []string
	}{
		messagePrefix: "nested error: root error",
		goReprPrefix:  "&fmt.wrapError{msg:\"nested error: root error\", err:",
		frames:        []string{"gosplit_test.errGoSplit2", "gosplit_test.errGoSplit1"},
	}

	err := errGoSplit2()
	detailed := fmt.Sprintf("%+v", err)
	goReprDetailed := fmt.Sprintf("%#+v", err)
	if ok := strings.HasPrefix(detailed, want.messagePrefix); !ok {
		t.Errorf("detailed = %#v, want HasPrefix(detailed, %#v)", detailed, want.messagePrefix)
	}
	if ok := strings.HasPrefix(goReprDetailed, want.goReprPrefix); !ok {
		t.Errorf("goReprDetailed = %#v, want HasPrefix(goReprDetailed, %#v)", goReprDetailed, want.goReprPrefix)
	}
	for _, frame := range want.frames {
		if strings.Count(detailed, frame) != 1 {
			t.Errorf("detailed = %#v, want Count(detailed, %#v) == 1", detailed, frame)
		}
		if strings.Count(goReprDetailed, frame) != 1 {
			t.Errorf("goReprDetailed = %#v, want Count(goReprDetailed, %#v) == 1", goReprDetailed, frame)
		}
	}
}

func TestErrGoSplitIs(t *testing.T) {
	t.Parallel()

	rootErr := errors.New("an error")
	err := gosplit.GoSplitErrorf("TestErrGoSplitIs: %w", rootErr)
	if ok := errors.Is(err, gosplit.ErrGoSplit); !ok {
		t.Errorf("errors.Is(gosplit.ErrGoSplit) should return true")
	}
	if ok := errors.Is(err, rootErr); !ok {
		t.Errorf("errors.Is(rootErr) should return true")
	}
}
