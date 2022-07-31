//go:build go1.16
// +build go1.16

package gobdd

import (
	"fmt"
	"io"
	"io/fs"
)

// WithFeaturesFS configures a filesystem and a pattern (regexp) where features can be found.
// An empty path defaults to "*.feature"
func WithFeaturesFS(fs fs.FS, pattern string) func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.featureSource = fsFeatureSource{
			fs:      fs,
			pattern: pattern,
		}
	}
}

type fsFeatureSource struct {
	fs      fs.FS
	pattern string
}

func (s fsFeatureSource) loadFeatures() ([]feature, error) {
	pattern := s.pattern
	if pattern == "" {
		pattern = "*.feature"
	}

	files, err := fs.Glob(s.fs, pattern)
	if err != nil {
		return nil, fmt.Errorf("loading features: %w", err)
	}

	features := make([]feature, 0, len(files))

	for _, f := range files {
		features = append(features, fsFeature{
			fs:   s.fs,
			file: f,
		})
	}

	return features, nil
}

type fsFeature struct {
	fs   fs.FS
	file string
}

func (f fsFeature) Open() (io.Reader, error) {
	file, err := f.fs.Open(string(f.file))
	if err != nil {
		return nil, fmt.Errorf("opening feature: %w", err)
	}

	return file, nil
}
