//go:build go1.16
// +build go1.16

package gobdd

import (
	"fmt"
	"io"
	"io/fs"
)

// WithFeaturesFS configures a filesystem and a path (glob pattern) where features can be found.
func WithFeaturesFS(fs fs.FS, patterns ...string) func(*SuiteOptions) {
	return func(options *SuiteOptions) {
		options.featureSource = fsFeatureSource{
			fs:       fs,
			patterns: patterns,
		}
	}
}

type fsFeatureSource struct {
	fs       fs.FS
	patterns []string
}

func (s fsFeatureSource) loadFeatures() ([]feature, error) {
	var allFiles []string
	for _, pattern := range s.patterns {
		files, err := fs.Glob(s.fs, pattern)
		if err != nil {
			return nil, fmt.Errorf("loading features: %w", err)
		}
		allFiles = append(allFiles, files...)
	}

	features := make([]feature, 0, len(allFiles))

	for _, f := range allFiles {
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
