package utils

// finds a string present in a slice
func Find(s []string, str string) []string {
	var matches []string
	if Contains(s, str) {
		matches = append(matches, str)
	}
	return matches
}

// contains checks if a string is present in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
