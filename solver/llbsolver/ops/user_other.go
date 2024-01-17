//go:build !linux && !windows
// +build !linux,!windows

package ops

import (
	"github.com/khulnasoft/devkit/snapshot"
	"github.com/khulnasoft/devkit/solver/pb"
	"github.com/khulnasoft/devkit/worker"
	"github.com/pkg/errors"
	copy "github.com/tonistiigi/fsutil/copy"
)

func getReadUserFn(worker worker.Worker) func(chopt *pb.ChownOpt, mu, mg snapshot.Mountable) (*copy.User, error) {
	return readUser
}

func readUser(chopt *pb.ChownOpt, mu, mg snapshot.Mountable) (*copy.User, error) {
	if chopt == nil {
		return nil, nil
	}
	return nil, errors.New("only implemented in linux and windows")
}
