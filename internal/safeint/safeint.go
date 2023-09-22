// Package safeint implements integer arithmetic checking integer overflow or division by zero.
package safeint

import (
	. "inaz2/GoSplit/internal/gerrors"
)

// MulInt64 returns the product x*y.
func MulInt64(x int64, y int64) (int64, Gerror) {
	z := x * y
	if y != 0 && z/y != x {
		return 0, SafeIntErrorf("%w: %#v * %#v -> %#v", ErrOverflow, x, y, z)
	}
	return z, nil
}

// PowInt64 returns the power b**k.
func PowInt64(b int64, k int64) (int64, Gerror) {
	var err Gerror

	if k < 0 {
		// b**k == 1 / (b**(-k))
		switch b {
		case 0:
			return 0, SafeIntErrorf("%w: 0 ** %#v", ErrDivisionByZero, k)
		case 1:
			return 1, nil
		case -1:
			if k&1 == 1 {
				return -1, nil
			} else {
				return 1, nil
			}
		default:
			return 0, nil
		}
	}

	result := int64(1)
	x := b
	for {
		if k&1 == 1 {
			result, err = MulInt64(result, x)
			if err != nil {
				return 0, err
			}
		}
		k >>= 1
		if k <= 0 {
			break
		}
		x, err = MulInt64(x, x)
		if err != nil {
			return 0, err
		}
	}
	return result, nil
}
