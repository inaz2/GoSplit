package gosplit

// safeMulInt64 return x*y with checking integer overflow
func safeMulInt64(x int64, y int64) (int64, bool) {
	z := x * y
	if y != 0 && z/y != x {
		// integer overflow occured
		return 0, false
	}
	return z, true
}

// safePowInt64 returns b**k with checking integer overflow
func safePowInt64(b int64, k int64) (int64, bool) {
	var ok bool

	if k < 0 {
		return 0, true
	}

	result := int64(1)
	x := b
	for {
		if k&1 == 1 {
			result, ok = safeMulInt64(result, x)
			if !ok {
				return 0, false
			}
		}
		k >>= 1
		if k <= 0 {
			break
		}
		x, ok = safeMulInt64(x, x)
		if !ok {
			return 0, false
		}
	}
	return result, true
}
