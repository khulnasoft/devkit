package main

import (
	"fmt"
	"os"

	_ "github.com/khulnasoft/devkit/client/connhelper/dockercontainer"
	_ "github.com/khulnasoft/devkit/client/connhelper/kubepod"
	_ "github.com/khulnasoft/devkit/client/connhelper/nerdctlcontainer"
	_ "github.com/khulnasoft/devkit/client/connhelper/podmancontainer"
	_ "github.com/khulnasoft/devkit/client/connhelper/ssh"
	bccommon "github.com/khulnasoft/devkit/cmd/buildctl/common"
	"github.com/khulnasoft/devkit/solver/errdefs"
	"github.com/khulnasoft/devkit/util/apicaps"
	"github.com/khulnasoft/devkit/util/appdefaults"
	"github.com/khulnasoft/devkit/util/profiler"
	"github.com/khulnasoft/devkit/util/stack"
	_ "github.com/khulnasoft/devkit/util/tracing/detect/delegated"
	_ "github.com/khulnasoft/devkit/util/tracing/detect/jaeger"
	_ "github.com/khulnasoft/devkit/util/tracing/env"
	"github.com/khulnasoft/devkit/version"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"go.opentelemetry.io/otel"
)

func init() {
	apicaps.ExportedProduct = "devkit"

	stack.SetVersionInfo(version.Version, version.Revision)

	// do not log tracing errors to stdio
	otel.SetErrorHandler(skipErrors{})
}

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(c.App.Name, version.Package, c.App.Version, version.Revision)
	}
	app := cli.NewApp()
	app.Name = "buildctl"
	app.Usage = "build utility"
	app.Version = version.Version

	defaultAddress := os.Getenv("BUILDKIT_HOST")
	if defaultAddress == "" {
		defaultAddress = appdefaults.Address
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in logs",
		},
		cli.StringFlag{
			Name:  "addr",
			Usage: "devkitd address",
			Value: defaultAddress,
		},
		// Add format flag to control log formatter
		cli.StringFlag{
			Name:  "log-format",
			Usage: "log formatter: json or text",
			Value: "text",
		},
		cli.StringFlag{
			Name:  "tlsservername",
			Usage: "devkitd server name for certificate validation",
			Value: "",
		},
		cli.StringFlag{
			Name:  "tlscacert",
			Usage: "CA certificate for validation",
			Value: "",
		},
		cli.StringFlag{
			Name:  "tlscert",
			Usage: "client certificate",
			Value: "",
		},
		cli.StringFlag{
			Name:  "tlskey",
			Usage: "client key",
			Value: "",
		},
		cli.StringFlag{
			Name:  "tlsdir",
			Usage: "directory containing CA certificate, client certificate, and client key",
			Value: "",
		},
		cli.IntFlag{
			Name:  "timeout",
			Usage: "timeout backend connection after value seconds",
			Value: 5,
		},
		cli.BoolFlag{
			Name:  "wait",
			Usage: "block RPCs until the connection becomes available",
		},
	}

	app.Commands = []cli.Command{
		diskUsageCommand,
		pruneCommand,
		pruneHistoriesCommand,
		buildCommand,
		debugCommand,
		dialStdioCommand,
	}

	var debugEnabled bool

	app.Before = func(context *cli.Context) error {
		debugEnabled = context.GlobalBool("debug")
		// Use Format flag to control log formatter
		logFormat := context.GlobalString("log-format")
		switch logFormat {
		case "json":
			logrus.SetFormatter(&logrus.JSONFormatter{})
		case "text", "":
			logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
		default:
			return errors.Errorf("unsupported log type %q", logFormat)
		}
		if debugEnabled {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}

	if err := bccommon.AttachAppContext(app); err != nil {
		handleErr(debugEnabled, err)
	}

	profiler.Attach(app)

	handleErr(debugEnabled, app.Run(os.Args))
}

func handleErr(debug bool, err error) {
	if err == nil {
		return
	}
	for _, s := range errdefs.Sources(err) {
		s.Print(os.Stderr)
	}
	if debug {
		fmt.Fprintf(os.Stderr, "error: %+v", stack.Formatter(err))
	} else {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	os.Exit(1)
}

type skipErrors struct{}

func (skipErrors) Handle(err error) {}
