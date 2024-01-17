package llbsolver

import (
	"context"

	"github.com/khulnasoft/devkit/solver/pb"
)

type SourcePolicyEvaluator interface {
	Evaluate(ctx context.Context, op *pb.Op) (bool, error)
}
