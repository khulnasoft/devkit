package http

import (
	"github.com/khulnasoft/devkit/solver/llbsolver/provenance"
	"github.com/khulnasoft/devkit/source"
	srctypes "github.com/khulnasoft/devkit/source/types"
	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

func NewHTTPIdentifier(str string, tls bool) (*HTTPIdentifier, error) {
	proto := "https://"
	if !tls {
		proto = "http://"
	}
	return &HTTPIdentifier{TLS: tls, URL: proto + str}, nil
}

type HTTPIdentifier struct {
	TLS      bool
	URL      string
	Checksum digest.Digest
	Filename string
	Perm     int
	UID      int
	GID      int
}

var _ source.Identifier = (*HTTPIdentifier)(nil)

func (id *HTTPIdentifier) Scheme() string {
	if id.TLS {
		return srctypes.HTTPSScheme
	}
	return srctypes.HTTPScheme
}

func (id *HTTPIdentifier) Capture(c *provenance.Capture, pin string) error {
	dgst, err := digest.Parse(pin)
	if err != nil {
		return errors.Wrapf(err, "failed to parse HTTP digest %s", pin)
	}
	c.AddHTTP(provenance.HTTPSource{
		URL:    id.URL,
		Digest: dgst,
	})
	return nil
}
