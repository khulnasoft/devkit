//go:build linux
// +build linux

package cniprovider

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	cni "github.com/containerd/go-cni"
	"github.com/khulnasoft/devkit/util/bklog"
	"github.com/khulnasoft/devkit/util/network"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

func NewBridge(opt Opt) (network.Provider, error) {
	cniOptions := []cni.Opt{cni.WithInterfacePrefix("eth")}
	bridgeBinName := "bridge"
	loopbackBinName := "loopback"
	hostLocalBinName := "host-local"
	firewallBinName := "firewall"
	var setup bool
	// binaries shipping with devkit
	for {
		var dirs []string

		bridgePath, err := exec.LookPath("devkit-cni-bridge")
		if err != nil {
			break
		}
		var bridgeDir string
		bridgeDir, bridgeBinName = filepath.Split(bridgePath)
		dirs = append(dirs, bridgeDir)

		loopbackPath, err := exec.LookPath("devkit-cni-loopback")
		if err != nil {
			break
		}
		var loopbackDir string
		loopbackDir, loopbackBinName = filepath.Split(loopbackPath)
		if loopbackDir != bridgeDir {
			dirs = append(dirs, loopbackDir)
		}

		hostLocalPath, err := exec.LookPath("devkit-cni-host-local")
		if err != nil {
			break
		}
		var hostLocalDir string
		hostLocalDir, hostLocalBinName = filepath.Split(hostLocalPath)
		if hostLocalDir != bridgeDir && hostLocalDir != loopbackDir {
			dirs = append(dirs, hostLocalDir)
		}

		firewallPath, err := exec.LookPath("devkit-cni-firewall")
		if err != nil {
			break
		}
		var firewallDir string
		firewallDir, firewallBinName = filepath.Split(firewallPath)
		if firewallDir != bridgeDir && firewallDir != loopbackDir && firewallDir != hostLocalDir {
			dirs = append(dirs, firewallDir)
		}

		cniOptions = append(cniOptions, cni.WithPluginDir(dirs))
		setup = true
		break //nolint: staticcheck
	}

	if !setup {
		fn := filepath.Join(opt.BinaryDir, "bridge")
		if _, err := os.Stat(fn); err != nil {
			return nil, errors.Wrapf(err, "failed to find CNI bridge %q or devkit-cni-bridge", fn)
		}

		cniOptions = append(cniOptions, cni.WithPluginDir([]string{opt.BinaryDir}))
	}

	cniOptions = append(cniOptions, cni.WithConfListBytes([]byte(fmt.Sprintf(`{
		"cniVersion": "1.0.0",
		"name": "devkit",
		"plugins": [
			{
				"type": "%s"
			},
			{
				"type": "%s",
				"bridge": "%s",
				"isDefaultGateway": true,
				"ipMasq": true,
				"ipam": {
				  "type": "%s",
				  "ranges": [
					[
					  { "subnet": "%s" }
					]
				  ]
				}
			  },
			  {
				"type": "%s",
				"ingressPolicy": "same-bridge"
			}
		]
		}`, loopbackBinName, bridgeBinName, opt.BridgeName, hostLocalBinName, opt.BridgeSubnet, firewallBinName))))

	unlock, err := initLock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	createBridge := true
	if _, err := bridgeByName(opt.BridgeName); err == nil {
		createBridge = false
	}

	cniHandle, err := cni.New(cniOptions...)
	if err != nil {
		return nil, err
	}
	cp := &cniProvider{
		CNI:  cniHandle,
		root: opt.Root,
	}

	if createBridge {
		cp.release = func() error {
			if err := removeBridge(opt.BridgeName); err != nil {
				bklog.L.Errorf("failed to remove bridge %q: %v", opt.BridgeName, err)
			}
			return nil
		}
	}

	cleanOldNamespaces(cp)

	cp.nsPool = &cniPool{targetSize: opt.PoolSize, provider: cp}
	if err := cp.initNetwork(false); err != nil {
		return nil, err
	}
	go cp.nsPool.fillPool(context.TODO())
	return cp, nil
}

func bridgeByName(name string) (*netlink.Bridge, error) {
	l, err := netlink.LinkByName(name)
	if err != nil {
		return nil, errors.Wrapf(err, "could not lookup %q", name)
	}
	br, ok := l.(*netlink.Bridge)
	if !ok {
		return nil, errors.Errorf("%q already exists but is not a bridge", name)
	}
	return br, nil
}

func removeBridge(name string) error {
	br, err := bridgeByName(name)
	if err != nil {
		return err
	}
	return netlink.LinkDel(br)
}
