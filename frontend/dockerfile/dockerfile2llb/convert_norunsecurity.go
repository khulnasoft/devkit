//go:build !dfrunsecurity
// +build !dfrunsecurity

package dockerfile2llb

import (
	"github.com/khulnasoft/devkit/client/llb"
	"github.com/khulnasoft/devkit/frontend/dockerfile/instructions"
)

func dispatchRunSecurity(c *instructions.RunCommand) (llb.RunOption, error) {
	return nil, nil
}
