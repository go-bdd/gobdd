package gobdd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextNilInGetError(t *testing.T) {
	ctx := NewContext()
	ctx.Set("err", nil)

	res, err := ctx.GetError("err")
	assert.NoError(t, err)
	assert.Nil(t, res)
}
