package layout_test

import (
	"gurms/internal/infra/logging/core/layout"
	"testing"
)

func TestFormatStructName(t *testing.T) {
	layout.FormatStructName("FrankWalterSteinmeier")
}
