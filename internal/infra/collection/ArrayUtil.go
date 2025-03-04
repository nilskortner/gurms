package collection

import "reflect"

func Concat(a, b []byte) []byte {
	var length = len(a)
	result := make([]byte, length+len(b))
	copy(a, result)
	copy(b, result[length:])
	return result
}

func GetSlice(value any) []any {
	// Check if the input value is already a slice.
	if outputSlice, ok := value.([]any); ok {
		return outputSlice
	}

	v := reflect.ValueOf(value)
	var length = v.Len()

	// If not a slice, initialize a new slice.
	outputSlice := make([]any, length)

	// Populate the slice with values from the input value.
	for i := 0; i < length; i++ {
		outputSlice[i] = v.Index(i).Interface()
	}

	return outputSlice
}
