//go:build freebsd || windows
// +build freebsd windows

package netproviders

import (
	"runtime"

	"github.com/khulnasoft/devkit/util/network"
	"github.com/khulnasoft/devkit/util/network/cniprovider"
	"github.com/pkg/errors"
)

func getBridgeProvider(opt cniprovider.Opt) (network.Provider, error) {
	return nil, errors.Errorf("bridge network is not supported on %s yet", runtime.GOOS)
}
