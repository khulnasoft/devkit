package proc

import (
	"context"

	"github.com/khulnasoft/devkit/client/llb"
	"github.com/khulnasoft/devkit/executor/resources"
	"github.com/khulnasoft/devkit/exporter/containerimage/exptypes"
	"github.com/khulnasoft/devkit/frontend"
	"github.com/khulnasoft/devkit/frontend/attestations/sbom"
	"github.com/khulnasoft/devkit/solver"
	"github.com/khulnasoft/devkit/solver/llbsolver"
	"github.com/khulnasoft/devkit/solver/result"
	"github.com/pkg/errors"
)

func SBOMProcessor(scannerRef string, useCache bool, resolveMode string) llbsolver.Processor {
	return func(ctx context.Context, res *llbsolver.Result, s *llbsolver.Solver, j *solver.Job, usage *resources.SysSampler) (*llbsolver.Result, error) {
		// skip sbom generation if we already have an sbom
		if sbom.HasSBOM(res.Result) {
			return res, nil
		}

		ps, err := exptypes.ParsePlatforms(res.Metadata)
		if err != nil {
			return nil, err
		}

		scanner, err := sbom.CreateSBOMScanner(ctx, s.Bridge(j), scannerRef, llb.ResolveImageConfigOpt{
			ResolveMode: resolveMode,
		})
		if err != nil {
			return nil, err
		}
		if scanner == nil {
			return res, nil
		}

		for _, p := range ps.Platforms {
			ref, ok := res.FindRef(p.ID)
			if !ok {
				return nil, errors.Errorf("could not find ref %s", p.ID)
			}
			if ref == nil {
				continue
			}

			defop, err := llb.NewDefinitionOp(ref.Definition())
			if err != nil {
				return nil, err
			}
			st := llb.NewState(defop)

			var opts []llb.ConstraintsOpt
			if !useCache {
				opts = append(opts, llb.IgnoreCache)
			}
			att, err := scanner(ctx, p.ID, st, nil, opts...)
			if err != nil {
				return nil, err
			}
			attSolve, err := result.ConvertAttestation(&att, func(st *llb.State) (solver.ResultProxy, error) {
				def, err := st.Marshal(ctx)
				if err != nil {
					return nil, err
				}

				r, err := s.Bridge(j).Solve(ctx, frontend.SolveRequest{
					Definition: def.ToPB(),
				}, j.SessionID)
				if err != nil {
					return nil, err
				}
				return r.Ref, nil
			})
			if err != nil {
				return nil, err
			}
			res.AddAttestation(p.ID, *attSolve)
		}
		return res, nil
	}
}
