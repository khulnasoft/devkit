package debug

import (
	"fmt"
	"os"
	"text/tabwriter"

	bccommon "github.com/khulnasoft/devkit/cmd/buildctl/common"
	"github.com/urfave/cli"
)

var InfoCommand = cli.Command{
	Name:   "info",
	Usage:  "display internal information",
	Action: info,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "Format the output using the given Go template, e.g, '{{json .}}'",
		},
	},
}

func info(clicontext *cli.Context) error {
	c, err := bccommon.ResolveClient(clicontext)
	if err != nil {
		return err
	}
	res, err := c.Info(bccommon.CommandContext(clicontext))
	if err != nil {
		return err
	}
	if format := clicontext.String("format"); format != "" {
		tmpl, err := bccommon.ParseTemplate(format)
		if err != nil {
			return err
		}
		if err := tmpl.Execute(clicontext.App.Writer, res); err != nil {
			return err
		}
		_, err = fmt.Fprintf(clicontext.App.Writer, "\n")
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	_, _ = fmt.Fprintf(w, "DevKit:\t%s %s %s\n", res.DevkitVersion.Package, res.DevkitVersion.Version, res.DevkitVersion.Revision)
	return w.Flush()
}
