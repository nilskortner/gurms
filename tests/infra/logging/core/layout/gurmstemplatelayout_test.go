package layout_test

import (
	"gurms/internal/infra/logging/core/layout"
	"testing"
)

func TestFormatStructName(t *testing.T) {
	result := string(layout.FormatStructName("Frank.Walter.Steinmeier"))

	t.Log(result)
}
