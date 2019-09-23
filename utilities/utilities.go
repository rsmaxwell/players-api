package utilities

// Equal tells whether a and b contain the same elements NOT in-order order
func Equal(x, y []string) bool {

	if x == nil {
		if y == nil {
			return true
		}
		return false
	} else if y == nil {
		return false
	}

	if len(x) != len(y) {
		return false
	}

	xMap := make(map[string]int)
	yMap := make(map[string]int)

	for _, xElem := range x {
		xMap[xElem]++
	}
	for _, yElem := range y {
		yMap[yElem]++
	}

	for xMapKey, xMapVal := range xMap {
		if yMap[xMapKey] != xMapVal {
			return false
		}
	}
	return true
}

// Equal2 tells whether a and b contain the same elements IN order
func Equal2(x, y []string) bool {

	if len(x) != len(y) {
		return false
	}
	for i, v := range x {
		if v != y[i] {
			return false
		}
	}
	return true
}
