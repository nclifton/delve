package stringutil

func Includes(a []string, b string) bool {
	found := false
	for _, c := range a {
		if c == b {
			found = true
		}
	}
	return found
}
