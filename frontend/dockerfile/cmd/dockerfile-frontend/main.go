package main

import (
	"flag"
	"fmt"
	"os"

	dockerfile "github.com/khulnasoft/devkit/frontend/dockerfile/builder"
	"github.com/khulnasoft/devkit/frontend/gateway/grpcclient"
	"github.com/khulnasoft/devkit/util/appcontext"
	"github.com/khulnasoft/devkit/util/bklog"
	"github.com/khulnasoft/devkit/util/stack"
)

func init() {
	stack.SetVersionInfo(Version, Revision)
}

func main() {
	var version bool
	flag.BoolVar(&version, "version", false, "show version")
	flag.Parse()

	if version {
		fmt.Printf("%s %s %s %s\n", os.Args[0], Package, Version, Revision)
		os.Exit(0)
	}

	if err := grpcclient.RunFromEnvironment(appcontext.Context(), dockerfile.Build); err != nil {
		bklog.L.Errorf("fatal error: %+v", err)
		panic(err)
	}
}
