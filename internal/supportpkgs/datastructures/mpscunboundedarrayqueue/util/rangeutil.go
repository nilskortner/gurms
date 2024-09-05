package util

import "fmt"

func CheckGreaterThanOrEqual(n, expected int, name string) (int, error) {
	if n < expected {
		return 0, fmt.Errorf("%s: %d (expected: >= %d)", name, n, expected)
	}

	return n, nil
}
