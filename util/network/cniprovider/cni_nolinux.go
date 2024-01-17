//go:build !linux
// +build !linux

package cniprovider

import (
	"github.com/khulnasoft/devkit/util/network"
)

func (ns *cniNS) sample() (*network.Sample, error) {
	return nil, nil
}
