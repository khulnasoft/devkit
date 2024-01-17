//go:build !linux

package resources

import "github.com/khulnasoft/devkit/executor/resources/types"

func newSysSampler() (*Sampler[*types.SysSample], error) {
	return nil, nil
}
