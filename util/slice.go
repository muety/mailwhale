package util

func ContainsString(needle string, haystack []string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
