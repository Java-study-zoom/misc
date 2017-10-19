package idutil

// Short returns the first 7 bytes of a string.
func Short(id string) string {
	if len(id) > 7 {
		return id[:7]
	}
	return id
}
