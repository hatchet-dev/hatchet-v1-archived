package modulerunner

import (
	"context"
)

// TODO: the runner should accept the same enviroment variables as the provisioner.
// There is no state here other than what is passed as env, nor is there DB access.
// type ModuleRunnerOpts struct {
// 	Provisioner provisioner.Provisioner
// 	Config      *server.Config
// }

// type ModuleRunner struct {
// 	prov provisioner.Provisioner
// 	conf *server.Config
// }

// func NewModuleRunner(opts *ModuleRunnerOpts) *ModuleRunner {
// 	return &ModuleRunner{opts.Provisioner, opts.Config}
// }

type RunInput struct {
	TeamID, ModuleID, ModuleRunID string
}

func Run(ctx context.Context, input RunInput) (string, error) {
	// construct provisioner options from config and input

	// call provisioner

	return "triggered", nil
}
