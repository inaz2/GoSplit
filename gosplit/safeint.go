package gosplit

// safeMulInt64 return x*y with checking integer overflow.
func safeMulInt64(x int64, y int64) (int64, error) {
	z := x * y
	if y != 0 && z/y != x {
		return 0, GoSplitErrorf("integer overflow occured: %#v * %#v -> %#v", x, y, z)
	}
	return z, nil
}

// safePowInt64 returns b**k with checking integer overflow.
func safePowInt64(b int64, k int64) (int64, error) {
	var err error

	if k < 0 {
		// b**k == 1 / (b**(-k))
		switch b {
		case 0:
			return 0, GoSplitErrorf("division by zero: 0 ** %#v", k)
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
			result, err = safeMulInt64(result, x)
			if err != nil {
				return 0, err
			}
		}
		k >>= 1
		if k <= 0 {
			break
		}
		x, err = safeMulInt64(x, x)
		if err != nil {
			return 0, err
		}
	}
	return result, nil
}
