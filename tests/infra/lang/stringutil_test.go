package lang_test

import (
	"gurms/internal/infra/lang"
	"testing"
)

func TestTokenize(t *testing.T) {
	str := lang.TokenizeToStringArray("frank.walter.steinmeier", ".")

	t.Log(str)
}
