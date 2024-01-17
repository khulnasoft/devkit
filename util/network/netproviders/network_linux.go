//go:build linux
// +build linux

package netproviders

import (
	"github.com/khulnasoft/devkit/util/network"
	"github.com/khulnasoft/devkit/util/network/cniprovider"
)

func getBridgeProvider(opt cniprovider.Opt) (network.Provider, error) {
	return cniprovider.NewBridge(opt)
}
