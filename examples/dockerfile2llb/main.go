package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"os"

	"github.com/khulnasoft/devkit/client/llb"
	"github.com/khulnasoft/devkit/client/llb/imagemetaresolver"
	"github.com/khulnasoft/devkit/frontend/dockerfile/dockerfile2llb"
	"github.com/khulnasoft/devkit/frontend/dockerui"
	"github.com/khulnasoft/devkit/solver/pb"
	"github.com/khulnasoft/devkit/util/appcontext"
	"github.com/sirupsen/logrus"
)

type buildOpt struct {
	target                 string
	partialImageConfigFile string
	partialMetadataFile    string
}

func main() {
	if err := xmain(); err != nil {
		logrus.Fatal(err)
	}
}

func xmain() error {
	var opt buildOpt
	flag.StringVar(&opt.target, "target", "", "target stage")
	flag.StringVar(&opt.partialImageConfigFile, "partial-image-config-file", "", "Output partial image config as a JSON file")
	flag.StringVar(&opt.partialMetadataFile, "partial-metadata-file", "", "Output partial metadata sa a JSON file")
	flag.Parse()

	df, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	caps := pb.Caps.CapSet(pb.Caps.All())

	state, img, _, err := dockerfile2llb.Dockerfile2LLB(appcontext.Context(), df, dockerfile2llb.ConvertOpt{
		MetaResolver: imagemetaresolver.Default(),
		LLBCaps:      &caps,
		Config: dockerui.Config{
			Target: opt.target,
		},
	})
	if err != nil {
		return err
	}

	dt, err := state.Marshal(context.TODO())
	if err != nil {
		return err
	}
	if err := llb.WriteTo(dt, os.Stdout); err != nil {
		return err
	}
	if opt.partialImageConfigFile != "" {
		if err := writeJSON(opt.partialImageConfigFile, img); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(f string, x interface{}) error {
	b, err := json.Marshal(x)
	if err != nil {
		return err
	}
	if err := os.RemoveAll(f); err != nil {
		return err
	}
	return os.WriteFile(f, b, 0o644)
}
