package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	"github.com/khulnasoft/devkit/client/llb"
	"github.com/khulnasoft/devkit/util/system"
)

type buildOpt struct {
	devkit       string
	containerd     string
	runc           string
	withContainerd bool
	installPrefix  string
}

func main() {
	var opt buildOpt
	flag.BoolVar(&opt.withContainerd, "with-containerd", true, "enable containerd worker")
	flag.StringVar(&opt.containerd, "containerd", "v1.7.2", "containerd version")
	flag.StringVar(&opt.runc, "runc", "v1.1.7", "runc version")
	flag.StringVar(&opt.devkit, "devkit", "master", "devkit version")
	flag.StringVar(&opt.installPrefix, "prefix", "/usr/local/bin", "path under which binaries should be installed")
	flag.Parse()

	bk := devkit(opt)
	out := bk
	dt, err := out.Marshal(context.TODO())
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
		Run(llb.Shlex("apk add --no-cache g++ linux-headers libseccomp-dev make")).Root()
}

func goRepo(s llb.State, repo string, src llb.State) func(ro ...llb.RunOption) llb.State {
	dir := "/go/src/" + repo
	return func(ro ...llb.RunOption) llb.State {
		es := s.Dir(dir).Run(ro...)
		es.AddMount(dir, src, llb.Readonly)
		return es.AddMount("/out", llb.Scratch())
	}
}

func runc(version string) llb.State {
	repo := "github.com/opencontainers/runc"
	src := llb.Git(repo, version)
	if version == "local" {
		src = llb.Local("runc-src")
	}
	return goRepo(goBuildBase(), repo, src)(
		llb.Shlex("go build -o /out/runc ./"),
	)
}

func containerd(version string) llb.State {
	repo := "github.com/containerd/containerd"
	src := llb.Git(repo, version, llb.KeepGitDir())
	if version == "local" {
		src = llb.Local("containerd-src")
	}
	return goRepo(
		goBuildBase().
			Run(llb.Shlex("apk add --no-cache btrfs-progs-dev")).Root(),
		repo, src)(
		llb.Shlex("go build -o /out/containerd ./cmd/containerd"),
	)
}

func devkit(opt buildOpt) llb.State {
	repo := "github.com/khulnasoft/devkit"
	src := llb.Git(repo, opt.devkit)
	if opt.devkit == "local" {
		src = llb.Local("devkit-src")
	}
	run := goRepo(goBuildBase(), repo, src)

	devkitd := prefixed(run(llb.Shlex("go build -o /out/devkitd ./cmd/devkitd")), opt.installPrefix)

	buildctl := prefixed(run(llb.Shlex("go build -o /out/buildctl ./cmd/buildctl")), opt.installPrefix)

	inputs := []llb.State{buildctl, devkitd, prefixed(runc(opt.runc), opt.installPrefix)}

	if opt.withContainerd {
		inputs = append(inputs, prefixed(containerd(opt.containerd), opt.installPrefix), devkitd)
	}
	return llb.Merge(inputs)
}

func prefixed(st llb.State, prefix string) llb.State {
	prefix = filepath.Clean(prefix)
	if prefix == "/" {
		return st
	}
	return llb.Scratch().File(llb.Copy(st, "/", prefix, &llb.CopyInfo{CreateDestPath: true}))
}