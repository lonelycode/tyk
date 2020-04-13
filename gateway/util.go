package gateway

// appendIfMissing appends the given new item to the given slice.
func appendIfMissing(slice []string, newSlice ...string) []string {
	for _, new := range newSlice {
		found := false
		for _, item := range slice {
			if item == new {
				continue
			}
			found = true
		}

		if !found {
			slice = append(slice, new)
		}
	}

	return slice
}

// contains checks whether the given slice contains the given item.
func contains(s []string, i string) bool {
	for _, a := range s {
		if a == i {
			return true
		}
	}
	return false
}

// greaterThanFloat64 checks whether first float64 value is bigger than second float64 value.
// -1 means infinite and the biggest value.
func greaterThanFloat64(first, second float64) bool {
	if first == -1 {
		return true
	}

	if second == -1 {
		return false
	}

	return first > second
}

// greaterThanInt64 checks whether first int64 value is bigger than second int64 value.
// -1 means infinite and the biggest value.
func greaterThanInt64(first, second int64) bool {
	if first == -1 {
		return true
	}

	if second == -1 {
		return false
	}

	return first > second
}
