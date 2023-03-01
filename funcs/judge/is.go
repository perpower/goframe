package judge

import "reflect"

// IsInt checks whether `value` is type of int.
func IsInt(value interface{}) bool {
	switch value.(type) {
	case int, *int, int8, *int8, int16, *int16, int32, *int32, int64, *int64:
		return true
	}
	return false
}

// IsUint checks whether `value` is type of uint.
func IsUint(value interface{}) bool {
	switch value.(type) {
	case uint, *uint, uint8, *uint8, uint16, *uint16, uint32, *uint32, uint64, *uint64:
		return true
	}
	return false
}

// IsFloat checks whether `value` is type of float.
func IsFloat(value interface{}) bool {
	switch value.(type) {
	case float32, *float32, float64, *float64:
		return true
	}
	return false
}

// IsSlice checks whether `value` is type of slice.
func IsSlice(value interface{}) bool {
	var (
		reflectValue = reflect.ValueOf(value)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		return true
	}
	return false
}

// IsMap checks whether `value` is type of map.
func IsMap(value interface{}) bool {
	var (
		reflectValue = reflect.ValueOf(value)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	switch reflectKind {
	case reflect.Map:
		return true
	}
	return false
}

// IsStruct checks whether `value` is type of struct.
func IsStruct(value interface{}) bool {
	var reflectType = reflect.TypeOf(value)
	if reflectType == nil {
		return false
	}
	var reflectKind = reflectType.Kind()
	for reflectKind == reflect.Ptr {
		reflectType = reflectType.Elem()
		reflectKind = reflectType.Kind()
	}
	switch reflectKind {
	case reflect.Struct:
		return true
	}
	return false
}

// IsNumeric checks whether the given string s is numeric.
// Note that float string like "123.456" is also numeric.
func IsNumeric(s string) bool {
	var (
		dotCount = 0
		length   = len(s)
	)
	if length == 0 {
		return false
	}
	for i := 0; i < length; i++ {
		if s[i] == '-' && i == 0 {
			continue
		}
		if s[i] == '.' {
			dotCount++
			if i > 0 && i < length-1 {
				continue
			} else {
				return false
			}
		}
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	if dotCount > 1 {
		return false
	}
	return true
}
