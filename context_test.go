package gobdd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContextNilInGetError(t *testing.T) {
	ctx := NewContext()
	ctx.Set("err", nil)

	res, err := ctx.GetError("err")
	require.NoError(t, err)
	require.NoError(t, res)
}
