package goutils

// ParseSpecialChars will parse any special chars to a printable version of it, just magic
func ParseSpecialChars(original string) string {
	var parsed []byte
	for _, c := range original {
		parsed = append(parsed, byte(c))
	}
	return string(parsed)
}
