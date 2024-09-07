package testutils_test

import (
	"testing"

	"github.com/fgognet/prompt/internal/testutils"
	"github.com/maxatome/go-testdeep/td"
)

func TestPromptsUtil(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		promptFunc := testutils.Prompts("foo", "bar")
		var i1 string
		var i2 string
		_, err := promptFunc(&i1)
		td.CmpNil(t, err)
		td.Cmp(t, i1, "foo")
		_, err = promptFunc()
		td.CmpNil(t, err)
		_, err = promptFunc(&i2)
		td.CmpNil(t, err)
		td.Cmp(t, i2, "bar")
		_, err = promptFunc()
		td.CmpNil(t, err)
		var i3 string
		_, err = promptFunc(&i3)
		td.CmpError(t, err, "not enough data")
	})
}
