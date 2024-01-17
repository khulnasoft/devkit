package frontend

import (
	"context"

	"github.com/khulnasoft/devkit/client/llb"
	gw "github.com/khulnasoft/devkit/frontend/gateway/client"
	"github.com/khulnasoft/devkit/session"
	"github.com/khulnasoft/devkit/solver"
	"github.com/khulnasoft/devkit/solver/pb"
	"github.com/khulnasoft/devkit/solver/result"
	digest "github.com/opencontainers/go-digest"
)

type Result = result.Result[solver.ResultProxy]

type Attestation = result.Attestation[solver.ResultProxy]

type Frontend interface {
	Solve(ctx context.Context, llb FrontendLLBBridge, opt map[string]string, inputs map[string]*pb.Definition, sid string, sm *session.Manager) (*Result, error)
}

type FrontendLLBBridge interface {
	Solve(ctx context.Context, req SolveRequest, sid string) (*Result, error)
	ResolveImageConfig(ctx context.Context, ref string, opt llb.ResolveImageConfigOpt) (string, digest.Digest, []byte, error)
	Warn(ctx context.Context, dgst digest.Digest, msg string, opts WarnOpts) error
}

type SolveRequest = gw.SolveRequest

type CacheOptionsEntry = gw.CacheOptionsEntry

type WarnOpts = gw.WarnOpts
