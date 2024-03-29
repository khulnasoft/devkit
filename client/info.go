package client

import (
	"context"

	controlapi "github.com/khulnasoft/devkit/api/services/control"
	apitypes "github.com/khulnasoft/devkit/api/types"
	"github.com/pkg/errors"
)

type Info struct {
	DevkitVersion DevkitVersion `json:"devkitVersion"`
}

type DevkitVersion struct {
	Package  string `json:"package"`
	Version  string `json:"version"`
	Revision string `json:"revision"`
}

func (c *Client) Info(ctx context.Context) (*Info, error) {
	res, err := c.ControlClient().Info(ctx, &controlapi.InfoRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to call info")
	}
	return &Info{
		DevkitVersion: fromAPIDevkitVersion(res.DevkitVersion),
	}, nil
}

func fromAPIDevkitVersion(in *apitypes.DevkitVersion) DevkitVersion {
	if in == nil {
		return DevkitVersion{}
	}
	return DevkitVersion{
		Package:  in.Package,
		Version:  in.Version,
		Revision: in.Revision,
	}
}
