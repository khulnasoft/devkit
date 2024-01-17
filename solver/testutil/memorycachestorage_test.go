package testutil

import (
	"testing"

	"github.com/khulnasoft/devkit/solver"
)

func TestMemoryCacheStorage(t *testing.T) {
	RunCacheStorageTests(t, solver.NewInMemoryCacheStorage)
}
