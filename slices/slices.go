package slices

// Last returns the value of the first index in s and true if s not empty,
// or nil and false if s empty.
func First[E any](s []E) (E, bool) {
	if len(s) == 0 {
		var zero E
		return zero, false
	}
	return s[0], true
}

// Last returns the value of the last index in s and true if s not empty,
// or nil and false if s empty.
func Last[E any](s []E) (E, bool) {
	if len(s) == 0 {
		var zero E
		return zero, false
	}
	return s[len(s)-1], true
}
