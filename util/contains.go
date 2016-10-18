package util

// Contains searches a slice and returns whether the slice contains the specified value.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
