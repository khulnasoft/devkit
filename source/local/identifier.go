package local

import (
	"github.com/khulnasoft/devkit/solver/llbsolver/provenance"
	"github.com/khulnasoft/devkit/source"
	srctypes "github.com/khulnasoft/devkit/source/types"
	"github.com/tonistiigi/fsutil"
)

type LocalIdentifier struct {
	Name            string
	SessionID       string
	IncludePatterns []string
	ExcludePatterns []string
	FollowPaths     []string
	SharedKeyHint   string
	Differ          fsutil.DiffType
}

func NewLocalIdentifier(str string) (*LocalIdentifier, error) {
	return &LocalIdentifier{Name: str}, nil
}

func (*LocalIdentifier) Scheme() string {
	return srctypes.LocalScheme
}

var _ source.Identifier = (*LocalIdentifier)(nil)

func (id *LocalIdentifier) Capture(c *provenance.Capture, pin string) error {
	c.AddLocal(provenance.LocalSource{
		Name: id.Name,
	})
	return nil
}
