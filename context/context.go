package context

import (
	gobdd "github.com/go-bdd/gobdd"
)

// Holds data from previously executed steps
// Deprecated: use gobdd.Context instead.
type Context = gobdd.Context

// Creates a new (empty) context struct
// Deprecated: use gobdd.NewContext instead.
func New() Context {
	return gobdd.NewContext()
}
