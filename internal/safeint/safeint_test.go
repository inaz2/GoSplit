package safeint_test

import (
	"inaz2/GoSplit/internal/safeint"

	"testing"
)

func TestMulInt64(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in        [2]int64
		want      int64
		expectErr bool
	}{
		"(-1)*(-1)":    {[2]int64{-1, -1}, 1, false},
		"(-1)*0":       {[2]int64{-1, 0}, 0, false},
		"(-1)*1":       {[2]int64{-1, 1}, -1, false},
		"0*(-1)":       {[2]int64{0, -1}, 0, false},
		"0*0":          {[2]int64{0, 0}, 0, false},
		"0*1":          {[2]int64{0, 1}, 0, false},
		"1*(-1)":       {[2]int64{1, -1}, -1, false},
		"1*0":          {[2]int64{1, 0}, 0, false},
		"1*1":          {[2]int64{1, 1}, 1, false},
		"(-4)*(1<<61)": {[2]int64{-4, 1 << 61}, (-4) * (1 << 61), false},
		"(-8)*(1<<61)": {[2]int64{-8, 1 << 61}, 0, true},
		"2*(1<<61)":    {[2]int64{2, 1 << 61}, 1 << 62, false},
		"4*(1<<61)":    {[2]int64{8, 1 << 61}, 0, true},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := safeint.MulInt64(tt.in[0], tt.in[1])
			if tt.expectErr && err == nil {
				t.Fatal("want err")
			}
			if !tt.expectErr && err != nil {
				t.Fatal("not want err:", err)
			}
			if tt.want != got {
				t.Errorf("MulInt64(%#v, %#v) = %#v, want %#v", tt.in[0], tt.in[1], got, tt.want)
			}
		})
	}
}

func TestPowInt64(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in        [2]int64
		want      int64
		expectErr bool
	}{
		"(-2)**(-1)":  {[2]int64{-2, -1}, 0, false},
		"(-1)**(-2)":  {[2]int64{-1, -2}, 1, false},
		"(-1)**(-1)":  {[2]int64{-1, -1}, -1, false},
		"(-1)**0":     {[2]int64{-1, 0}, 1, false},
		"(-1)**1":     {[2]int64{-1, 1}, -1, false},
		"0**(-1)":     {[2]int64{0, -1}, 0, true},
		"0**0":        {[2]int64{0, 0}, 1, false},
		"0**1":        {[2]int64{0, 1}, 0, false},
		"1**(-2)":     {[2]int64{1, -2}, 1, false},
		"1**(-1)":     {[2]int64{1, -1}, 1, false},
		"1**0":        {[2]int64{1, 0}, 1, false},
		"1**1":        {[2]int64{1, 1}, 1, false},
		"2**(-1)":     {[2]int64{2, -1}, 0, false},
		"(-2)**(-64)": {[2]int64{-2, -64}, 0, false},
		"(-2)**63":    {[2]int64{-2, 63}, (-2) * (1 << 62), false},
		"(-2)**64":    {[2]int64{-2, 64}, 0, true},
		"2**(-63)":    {[2]int64{2, -63}, 0, false},
		"2**62":       {[2]int64{2, 62}, 1 << 62, false},
		"2**63":       {[2]int64{2, 63}, 0, true},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := safeint.PowInt64(tt.in[0], tt.in[1])
			if tt.expectErr && err == nil {
				t.Fatal("want err")
			}
			if !tt.expectErr && err != nil {
				t.Fatal("not want err:", err)
			}
			if tt.want != got {
				t.Errorf("PowInt64(%#v, %#v) = %#v, want %#v", tt.in[0], tt.in[1], got, tt.want)
			}
		})
	}
}
