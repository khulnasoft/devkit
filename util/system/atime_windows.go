package system

import (
	iofs "io/fs"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

func Atime(st iofs.FileInfo) (time.Time, error) {
	stSys, ok := st.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return time.Time{}, errors.Errorf("expected st.Sys() to be *syscall.Win32FileAttributeData, got %T", st.Sys())
	}
	// ref: https://github.com/golang/go/blob/go1.19.2/src/os/types_windows.go#L230
	return time.Unix(0, stSys.LastAccessTime.Nanoseconds()), nil
}
