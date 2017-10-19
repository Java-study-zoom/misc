package strutil

// MakeSet converts a list of strings to a set of strings.
func MakeSet(lst []string) map[string]bool {
	ret := make(map[string]bool)
	for _, s := range lst {
		ret[s] = true
	}
	return ret
}
