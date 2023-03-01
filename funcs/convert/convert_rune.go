package convert

// Rune converts `any` to rune.
func Rune(any interface{}) rune {
	if v, ok := any.(rune); ok {
		return v
	}
	return Int32(any)
}

// Runes converts `any` to []rune.
func Runes(any interface{}) []rune {
	if v, ok := any.([]rune); ok {
		return v
	}
	return []rune(String(any))
}
