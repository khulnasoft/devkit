package dockerui

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/khulnasoft/devkit/frontend/gateway/client"
	"github.com/khulnasoft/devkit/frontend/subrequests"
	"github.com/khulnasoft/devkit/frontend/subrequests/outline"
	"github.com/khulnasoft/devkit/frontend/subrequests/targets"
	"github.com/khulnasoft/devkit/solver/errdefs"
)

const (
	keyRequestID = "requestid"
)

type RequestHandler struct {
	Outline     func(context.Context) (*outline.Outline, error)
	ListTargets func(context.Context) (*targets.List, error)
	AllowOther  bool
}

func (bc *Client) HandleSubrequest(ctx context.Context, h RequestHandler) (*client.Result, bool, error) {
	req, ok := bc.bopts.Opts[keyRequestID]
	if !ok {
		return nil, false, nil
	}
	switch req {
	case subrequests.RequestSubrequestsDescribe:
		res, err := describe(h)
		return res, true, err
	case outline.SubrequestsOutlineDefinition.Name:
		if f := h.Outline; f != nil {
			o, err := f(ctx)
			if err != nil {
				return nil, false, err
			}
			if o == nil {
				return nil, true, nil
			}
			res, err := o.ToResult()
			return res, true, err
		}
	case targets.SubrequestsTargetsDefinition.Name:
		if f := h.ListTargets; f != nil {
			targets, err := f(ctx)
			if err != nil {
				return nil, false, err
			}
			if targets == nil {
				return nil, true, nil
			}
			res, err := targets.ToResult()
			return res, true, err
		}
	}
	if h.AllowOther {
		return nil, false, nil
	}
	return nil, false, errdefs.NewUnsupportedSubrequestError(req)
}

func describe(h RequestHandler) (*client.Result, error) {
	all := []subrequests.Request{}
	if h.Outline != nil {
		all = append(all, outline.SubrequestsOutlineDefinition)
	}
	if h.ListTargets != nil {
		all = append(all, targets.SubrequestsTargetsDefinition)
	}
	all = append(all, subrequests.SubrequestsDescribeDefinition)
	dt, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(nil)
	if err := subrequests.PrintDescribe(dt, b); err != nil {
		return nil, err
	}

	res := client.NewResult()
	res.Metadata = map[string][]byte{
		"result.json": dt,
		"result.txt":  b.Bytes(),
		"version":     []byte(subrequests.SubrequestsDescribeDefinition.Version),
	}
	return res, nil
}
