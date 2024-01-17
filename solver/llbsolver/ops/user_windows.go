package ops

import (
	"context"

	"github.com/docker/docker/pkg/idtools"
	"github.com/khulnasoft/devkit/snapshot"
	"github.com/khulnasoft/devkit/solver/pb"
	"github.com/khulnasoft/devkit/util/windows"
	"github.com/khulnasoft/devkit/worker"
	"github.com/pkg/errors"
	copy "github.com/tonistiigi/fsutil/copy"
)

func getReadUserFn(worker worker.Worker) func(chopt *pb.ChownOpt, mu, mg snapshot.Mountable) (*copy.User, error) {
	return func(chopt *pb.ChownOpt, mu, mg snapshot.Mountable) (*copy.User, error) {
		return readUser(chopt, mu, mg, worker)
	}
}

func readUser(chopt *pb.ChownOpt, mu, mg snapshot.Mountable, worker worker.Worker) (*copy.User, error) {
	if chopt == nil {
		return nil, nil
	}

	if chopt.User != nil {
		switch u := chopt.User.User.(type) {
		case *pb.UserOpt_ByName:
			if mu == nil {
				return nil, errors.Errorf("invalid missing user mount")
			}

			rootMounts, release, err := mu.Mount()
			if err != nil {
				return nil, err
			}
			defer release()
			ident, err := windows.ResolveUsernameToSID(context.Background(), worker.Executor(), rootMounts, u.ByName.Name)
			if err != nil {
				return nil, err
			}
			return &copy.User{SID: ident.SID}, nil
		default:
			return &copy.User{SID: idtools.ContainerAdministratorSidString}, nil
		}
	}
	return &copy.User{SID: idtools.ContainerAdministratorSidString}, nil
}
