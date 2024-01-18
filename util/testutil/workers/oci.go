package workers

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/khulnasoft/devkit/util/testutil/integration"
	"github.com/pkg/errors"
)

// InitOCIWorker registers an integration test worker, which enables the --oci-worker
// flag in the test devkitd instance and disables the --containerd-worker flag. This
// integration test worker is not supported on Windows.
func InitOCIWorker() {
	// calling platform specific
	initOCIWorker()
}

type OCI struct {
	ID          string
	UID         int
	GID         int
	Snapshotter string
}

func (s *OCI) Name() string {
	return s.ID
}

func (s *OCI) Rootless() bool {
	return s.UID != 0
}

func (s *OCI) New(ctx context.Context, cfg *integration.BackendConfig) (integration.Backend, func() error, error) {
	if err := integration.LookupBinary("devkitd"); err != nil {
		return nil, nil, err
	}
	if err := requireRoot(); err != nil {
		return nil, nil, err
	}
	// Include use of --oci-worker-labels to trigger https://github.com/khulnasoft/devkit/pull/603
	devkitdArgs := []string{"devkitd", "--oci-worker=true", "--containerd-worker=false", "--oci-worker-gc=false", "--oci-worker-labels=org.khulnasoftproject.devkit.worker.sandbox=true"}

	if s.Snapshotter != "" {
		devkitdArgs = append(devkitdArgs,
			fmt.Sprintf("--oci-worker-snapshotter=%s", s.Snapshotter))
	}

	if s.UID != 0 {
		if s.GID == 0 {
			return nil, nil, errors.Errorf("unsupported id pair: uid=%d, gid=%d", s.UID, s.GID)
		}
		// TODO: make sure the user exists and subuid/subgid are configured.
		devkitdArgs = append([]string{"sudo", "-u", fmt.Sprintf("#%d", s.UID), "-i", "--", "exec", "rootlesskit"}, devkitdArgs...)
	}

	var extraEnv []string
	if runtime.GOOS != "windows" && s.Snapshotter != "native" {
		extraEnv = append(extraEnv, "DEVKIT_DEBUG_FORCE_OVERLAY_DIFF=true")
	}
	devkitdSock, stop, err := runDevkitd(ctx, cfg, devkitdArgs, cfg.Logs, s.UID, s.GID, extraEnv)
	if err != nil {
		integration.PrintLogs(cfg.Logs, log.Println)
		return nil, nil, err
	}

	return backend{
		address:     devkitdSock,
		rootless:    s.UID != 0,
		snapshotter: s.Snapshotter,
	}, stop, nil
}

func (s *OCI) Close() error {
	return nil
}
