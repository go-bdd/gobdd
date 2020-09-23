package gobdd

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContext_GetAs_NoPointerType(t *testing.T) {
	s := []string{"one", "two"}
	ctx := NewContext()
	ctx.Set("s", s)
	var res string
	err := ctx.GetAs("s", res)
	require.Error(t, err)
}

func TestContext_GetAs_WithSlice(t *testing.T) {
	s := []string{"one", "two"}
	ctx := NewContext()
	ctx.Set("s", s)
	var res []string
	err := ctx.GetAs("s", &res)
	require.NoError(t, err)
	require.Equal(t, s, res)
}

func TestContext_GetAs_WithMap(t *testing.T) {
	m := map[string]string{}
	m["key"] = "value"
	ctx := NewContext()
	ctx.Set("map", m)
	var res map[string]string
	err := ctx.GetAs("map", &res)
	require.NoError(t, err)
	require.Equal(t, m, res)
}