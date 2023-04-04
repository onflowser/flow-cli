/*
 * Flow CLI
 *
 * Copyright 2019 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/onflow/flow-cli/internal/util"

	"github.com/onflow/flow-cli/internal/command"
	"github.com/onflow/flow-cli/pkg/flowkit"
	"github.com/onflow/flow-cli/pkg/flowkit/config"
	"github.com/onflow/flow-cli/pkg/flowkit/output"
)

type flagsAddDeployment struct {
	Network   string   `flag:"network" info:"Network name used for deployment"`
	Account   string   `flag:"account" info:"Account name used for deployment"`
	Contracts []string `flag:"contract" info:"Name of the contract to be deployed"`
}

var addDeploymentFlags = flagsAddDeployment{}

var AddDeploymentCommand = &command.Command{
	Cmd: &cobra.Command{
		Use:     "deployment",
		Short:   "Add deployment to configuration",
		Example: "flow config add deployment",
		Args:    cobra.NoArgs,
	},
	Flags: &addDeploymentFlags,
	RunS:  addDeployment,
}

func addDeployment(
	_ []string,
	globalFlags command.GlobalFlags,
	_ output.Logger,
	_ flowkit.Services,
	state *flowkit.State,
) (command.Result, error) {
	raw, flagsProvided, err := flagsToDeploymentData(addDeploymentFlags)
	if err != nil {
		return nil, err
	}

	if !flagsProvided {
		raw = util.NewDeploymentPrompt(*state.Networks(), state.Config().Accounts, *state.Contracts())
	}

	contracts := make([]config.ContractDeployment, len(raw.Contracts))
	for i, c := range raw.Contracts {
		contracts[i] = config.ContractDeployment{Name: c}
	}

	state.Deployments().AddOrUpdate(config.Deployment{
		Network:   raw.Network,
		Account:   raw.Account,
		Contracts: contracts,
	})

	err = state.SaveEdited(globalFlags.ConfigPaths)
	if err != nil {
		return nil, err
	}

	return &Result{
		result: "Deployment added to the configuration.\nYou can deploy using 'flow project deploy' command",
	}, nil
}

func flagsToDeploymentData(flags flagsAddDeployment) (*util.DeploymentData, bool, error) {
	if flags.Network == "" && flags.Account == "" && len(flags.Contracts) == 0 {
		return nil, false, nil
	}

	if flags.Network == "" {
		return nil, true, fmt.Errorf("network name must be provided")
	} else if flags.Account == "" {
		return nil, true, fmt.Errorf("account name must be provided")
	} else if len(flags.Contracts) == 0 {
		return nil, true, fmt.Errorf("at least one contract name must be provided")
	}

	return &util.DeploymentData{
		Network:   flags.Network,
		Account:   flags.Account,
		Contracts: flags.Contracts,
	}, true, nil
}
