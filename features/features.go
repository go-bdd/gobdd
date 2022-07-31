//go:build go1.16
// +build go1.16

package features

import (
	"embed"
	"io/fs"
)

//go:embed *.feature
var features embed.FS

// Features returns a filesystem that contains feature files.
func Features() fs.FS {
  return features
}
