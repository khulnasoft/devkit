package main

import (
	"context"
	"os"

	"github.com/khulnasoft/devkit/client/llb"
	"github.com/khulnasoft/devkit/client/llb/llbbuild"
	"github.com/khulnasoft/devkit/util/system"
)

const url = "https://gist.githubusercontent.com/tonistiigi/03b4049f8cc3de059bd2a1a1d8643714/raw/b5960995d570d8c6d94db527e805edc6d5854268/buildprs.go"

func main() {
	build := goBuildBase().
		Run(llb.Shlex("apk add --no-cache curl")).
		Run(llb.Shlexf("curl -o /buildprs.go \"%s\"", url))

	devkitRepo := "github.com/khulnasoft/devkit"

	build = build.Run(llb.Shlex("sh -c \"go run /buildprs.go > /out/devkit.llb.definition\""))
	build.AddMount("/go/src/"+devkitRepo, llb.Git(devkitRepo, "master"))
	pb := build.AddMount("/out", llb.Scratch())

	built := pb.With(llbbuild.Build())

	dt, err := llb.Image("docker.io/library/alpine:latest").Run(llb.Shlex("ls -l /out"), llb.AddMount("/out", built, llb.Readonly)).Marshal(context.TODO(), llb.LinuxAmd64)
	if err != nil {
		panic(err)
	}
	llb.WriteTo(dt, os.Stdout)
}

func goBuildBase() llb.State {
	goAlpine := llb.Image("docker.io/library/golang:1.21-alpine")
	return goAlpine.
		AddEnv("PATH", "/usr/local/go/bin:"+system.DefaultPathEnvUnix).
		AddEnv("GOPATH", "/go").
		Run(llb.Shlex("apk add --no-cache g++ linux-headers make")).Root()
}
