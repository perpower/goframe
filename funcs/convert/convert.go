// 类型转换函数
package convert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	// Empty strings.
	emptyStringMap = map[string]struct{}{
		"":      {},
		"0":     {},
		"no":    {},
		"off":   {},
		"false": {},
	}
)

// String converts `any` to string.
// It's most commonly used converting function.
func String(any interface{}) string {
	if any == nil {
		return ""
	}
	switch value := any.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	case time.Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *time.Time:
		if value == nil {
			return ""
		}
		return value.String()
	default:
		// Empty checks.
		if value == nil {
			return ""
		}
		if f, ok := value.(iString); ok {
			// If the variable implements the String() interface,
			// then use that interface to perform the conversion
			return f.String()
		}
		if f, ok := value.(iError); ok {
			// If the variable implements the Error() interface,
			// then use that interface to perform the conversion
			return f.Error()
		}
		// Reflect checks.
		var (
			rv   = reflect.ValueOf(value)
			kind = rv.Kind()
		)
		switch kind {
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return ""
			}
		case reflect.String:
			return rv.String()
		}
		if kind == reflect.Ptr {
			return String(rv.Elem().Interface())
		}
		// Finally, we use json.Marshal to convert.
		if jsonContent, err := json.Marshal(value); err != nil {
			return fmt.Sprint(value)
		} else {
			return string(jsonContent)
		}
	}
}

// Bool converts `any` to bool.
// It returns false if `any` is: false, "", 0, "false", "off", "no", empty slice/map.
func Bool(any interface{}) bool {
	if any == nil {
		return false
	}
	switch value := any.(type) {
	case bool:
		return value
	case []byte:
		if _, ok := emptyStringMap[strings.ToLower(string(value))]; ok {
			return false
		}
		return true
	case string:
		if _, ok := emptyStringMap[strings.ToLower(value)]; ok {
			return false
		}
		return true
	default:
		if f, ok := value.(iBool); ok {
			return f.Bool()
		}
		rv := reflect.ValueOf(any)
		switch rv.Kind() {
		case reflect.Ptr:
			return !rv.IsNil()
		case reflect.Map:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			return rv.Len() != 0
		case reflect.Struct:
			return true
		default:
			s := strings.ToLower(String(any))
			if _, ok := emptyStringMap[s]; ok {
				return false
			}
			return true
		}
	}
}
