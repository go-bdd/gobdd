//go:build go1.16
// +build go1.16

package gobdd

import (
	"regexp"
	"testing"

	"github.com/go-bdd/gobdd/features"
)

func TestWithFeaturesFS(t *testing.T) {
	suite := NewSuite(t, WithFeaturesFS(features.Features(), "example.feature"))
	compiled := regexp.MustCompile(`I add (\d+) and (\d+)`)
	suite.AddRegexStep(compiled, add)
	compiled = regexp.MustCompile(`the result should equal (\d+)`)
	suite.AddRegexStep(compiled, check)

	suite.Run()
}
