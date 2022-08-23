//go:build go1.16
// +build go1.16

package gobdd

import (
	"fmt"
	"io"
	"io/fs"
)

// WithFeaturesFS configures a filesystem where features can be found and loads all files mathing the "*.feature" pattern.
func WithFeaturesFS(fs fs.FS) func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.featureSource = fsFeatureSource{
			fs:      fs,
			pattern: "*.feature",
		}
	}
}

// WithFeaturesFS configures a filesystem and a pattern (regexp) where features can be found.
// An empty path defaults to "*.feature"
func WithFeaturesFSPattern(fs fs.FS, pattern string) func(*SuiteOptions) {
	if pattern == "" {
		pattern = "*.feature"
	}

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
	files, err := fs.Glob(s.fs, s.pattern)
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
	file, err := f.fs.Open(f.file)
	if err != nil {
		return nil, fmt.Errorf("opening feature: %w", err)
	}

	return file, nil
}
