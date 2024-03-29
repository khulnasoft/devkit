package subrequests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/khulnasoft/devkit/frontend/gateway/client"
	gwpb "github.com/khulnasoft/devkit/frontend/gateway/pb"
	"github.com/khulnasoft/devkit/solver/errdefs"
	"github.com/pkg/errors"
)

const RequestSubrequestsDescribe = "frontend.subrequests.describe"

var SubrequestsDescribeDefinition = Request{
	Name:        RequestSubrequestsDescribe,
	Version:     "1.0.0",
	Type:        TypeRPC,
	Description: "List available subrequest types",
	Metadata: []Named{
		{Name: "result.json"},
		{Name: "result.txt"},
	},
}

func Describe(ctx context.Context, c client.Client) ([]Request, error) {
	gwcaps := c.BuildOpts().Caps

	if err := (&gwcaps).Supports(gwpb.CapFrontendCaps); err != nil {
		return nil, errdefs.NewUnsupportedSubrequestError(RequestSubrequestsDescribe)
	}

	res, err := c.Solve(ctx, client.SolveRequest{
		FrontendOpt: map[string]string{
			"requestid":     RequestSubrequestsDescribe,
			"frontend.caps": "moby.devkit.frontend.subrequests",
		},
		Frontend: "dockerfile.v0",
	})
	if err != nil {
		var reqErr *errdefs.UnsupportedSubrequestError
		if errors.As(err, &reqErr) {
			return nil, err
		}
		var capErr *errdefs.UnsupportedFrontendCapError
		if errors.As(err, &capErr) {
			return nil, errdefs.NewUnsupportedSubrequestError(RequestSubrequestsDescribe)
		}
		return nil, err
	}

	dt, ok := res.Metadata["result.json"]
	if !ok {
		return nil, errors.Errorf("no result.json metadata in response")
	}

	var reqs []Request
	if err := json.Unmarshal(dt, &reqs); err != nil {
		return nil, errors.Wrap(err, "failed to parse describe result")
	}
	return reqs, nil
}

func PrintDescribe(dt []byte, w io.Writer) error {
	var d []Request
	if err := json.Unmarshal(dt, &d); err != nil {
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
	fmt.Fprintf(tw, "NAME\tVERSION\tDESCRIPTION\n")

	for _, r := range d {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", strings.TrimPrefix(r.Name, "frontend."), r.Version, r.Description)
	}
	return tw.Flush()
}
