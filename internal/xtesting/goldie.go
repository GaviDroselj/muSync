package xtesting

import (
	"testing"

	"github.com/sebdah/goldie/v2"
)

func NewGoldie(t *testing.T) *goldie.Goldie {
	return goldie.New(
		t,
		goldie.WithFixtureDir("testdata"),
		goldie.WithNameSuffix(".golden.json"),
		goldie.WithDiffEngine(goldie.ColoredDiff),
		goldie.WithTestNameForDir(true),
		goldie.WithSubTestNameForDir(true),
	)
}
