package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilInGetError(t *testing.T) {
	ctx := New()
	ctx.Set("err", nil)

	res, err := ctx.GetError("err")
	assert.NoError(t, err)
	assert.Nil(t, res)
}
