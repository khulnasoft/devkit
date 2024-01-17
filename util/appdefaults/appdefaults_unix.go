//go:build !windows
// +build !windows

package appdefaults

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	Address              = "unix:///run/devkit/devkitd.sock"
	Root                 = "/var/lib/devkit"
	ConfigDir            = "/etc/devkit"
	DefaultCNIBinDir     = "/opt/cni/bin"
	DefaultCNIConfigPath = "/etc/devkit/cni.json"
)

var (
	UserCNIConfigPath = filepath.Join(UserConfigDir(), "cni.json")
)

// UserAddress typically returns /run/user/$UID/devkit/devkitd.sock
func UserAddress() string {
	//  pam_systemd sets XDG_RUNTIME_DIR but not other dirs.
	xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if xdgRuntimeDir != "" {
		dirs := strings.Split(xdgRuntimeDir, ":")
		return "unix://" + filepath.Join(dirs[0], "devkit", "devkitd.sock")
	}
	return Address
}

// EnsureUserAddressDir sets sticky bit on XDG_RUNTIME_DIR if XDG_RUNTIME_DIR is set.
// See https://github.com/opencontainers/runc/issues/1694
func EnsureUserAddressDir() error {
	xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if xdgRuntimeDir != "" {
		dirs := strings.Split(xdgRuntimeDir, ":")
		dir := filepath.Join(dirs[0], "devkit")
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
		return os.Chmod(dir, 0700|os.ModeSticky)
	}
	return nil
}

// UserRoot typically returns /home/$USER/.local/share/devkit
func UserRoot() string {
	//  pam_systemd sets XDG_RUNTIME_DIR but not other dirs.
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		dirs := strings.Split(xdgDataHome, ":")
		return filepath.Join(dirs[0], "devkit")
	}
	home := os.Getenv("HOME")
	if home != "" {
		return filepath.Join(home, ".local", "share", "devkit")
	}
	return Root
}

// UserConfigDir returns dir for storing config. /home/$USER/.config/devkit/
func UserConfigDir() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "devkit")
	}
	home := os.Getenv("HOME")
	if home != "" {
		return filepath.Join(home, ".config", "devkit")
	}
	return ConfigDir
}

func TraceSocketPath(inUserNS bool) string {
	if inUserNS {
		if xrd := os.Getenv("XDG_RUNTIME_DIR"); xrd != "" {
			dirs := strings.Split(xrd, ":")
			return filepath.Join(dirs[0], "devkit", "otel-grpc.sock")
		}
	}
	return "/run/devkit/otel-grpc.sock"
}
