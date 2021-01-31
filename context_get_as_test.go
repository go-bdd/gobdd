package gobdd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext_GetAs_NoPointerType(t *testing.T) {
	s := []string{"one", "two"}
	ctx := NewContext()
	ctx.Set("s", s)
	res := ""
	err := ctx.GetAs("s", res)
	require.Error(t, err)
}

func TestContext_GetAs_WithSlice(t *testing.T) {
	s := []string{"one", "two"}
	ctx := NewContext()
	ctx.Set("s", s)
	res := map[string]string{}
	err := ctx.GetAs("s", &res)
	require.NoError(t, err)
	require.Equal(t, s, res)
}

func TestContext_GetAs_WithMap(t *testing.T) {
	m := map[string]string{}
	m["key"] = "value"
	ctx := NewContext()
	ctx.Set("map", m)
	res := map[string]string{}
	err := ctx.GetAs("map", &res)
	require.NoError(t, err)
	require.Equal(t, m, res)
}
