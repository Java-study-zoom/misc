package syntax

func makeMap(lst []string) map[string]struct{} {
	ret := make(map[string]struct{})
	for _, s := range lst {
		ret[s] = struct{}{}
	}
	return ret
}
