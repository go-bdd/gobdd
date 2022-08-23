//go:build go1.16
// +build go1.16

package gobdd

import (
	"embed"
)

//go:embed features/*.feature
var featuresFS embed.FS
