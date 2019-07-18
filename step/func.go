package step

import "github.com/go-bdd/gobdd/context"

// Every step function have to be compatible with this type
type Func func(ctx context.Context) error
