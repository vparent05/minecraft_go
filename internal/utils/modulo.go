package utils

// Mod returns a modulo b and
// works with negative numbers. I.e. Mod(-1, 8) returns 7
func Mod(a, b int) int {
	return (a%b + b) % b
}
