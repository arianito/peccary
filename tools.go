package peccary

import "strings"

func isLower(ch byte) bool {
	return ch >= 'a' && ch <= 'z'
}
func toLower(ch byte) byte {
	if ch >= 'A' && ch <= 'Z' {
		return ch + 32
	}
	return ch
}
func camelCase(s string) string {
	s = strings.TrimSpace(s)
	if !isLower(s[0]) {
		return string(toLower(s[0])) + s[1:]
	}
	return s
}
