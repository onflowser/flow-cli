/*
 * Flow CLI
 *
 * Copyright 2019-2021 Dapper Labs, Inc.
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

package project

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/onflow/flow-cli/pkg/flow/util"

	"github.com/onflow/flow-cli/internal/command"

	"github.com/onflow/flow-cli/pkg/flow"

	"github.com/onflow/flow-cli/pkg/flow/services"
	"github.com/spf13/cobra"
)

type flagsInit struct {
	ServicePrivateKey  string `flag:"service-priv-key" info:"Service account private key"`
	ServiceKeySigAlgo  string `default:"ECDSA_P256" flag:"service-sig-algo" info:"Service account key signature algorithm"`
	ServiceKeyHashAlgo string `default:"SHA3_256" flag:"service-hash-algo" info:"Service account key hash algorithm"`
	Reset              bool   `default:"false" flag:"reset" info:"Reset flow.json config file"`
}

var initFlag = flagsInit{}

var InitCommand = &command.Command{
	Cmd: &cobra.Command{
		Use:   "init",
		Short: "Initialize a new account profile",
	},
	Flags: &initFlag,
	Run: func(
		cmd *cobra.Command,
		args []string,
		services *services.Services,
	) (command.Result, error) {
		project, err := services.Project.Init(
			initFlag.Reset,
			initFlag.ServiceKeySigAlgo,
			initFlag.ServiceKeyHashAlgo,
			initFlag.ServicePrivateKey,
		)
		if err != nil {
			return nil, err
		}

		return &InitResult{project}, nil
	},
}

// InitResult result structure
type InitResult struct {
	*flow.Project
}

// JSON convert result to JSON
func (r *InitResult) JSON() interface{} {
	return r
}

// String convert result to string
func (r *InitResult) String() string {
	var b bytes.Buffer
	writer := tabwriter.NewWriter(&b, 0, 8, 1, '\t', tabwriter.AlignRight)
	account, _ := r.Project.EmulatorServiceAccount()

	fmt.Fprintf(writer, "Configuration initialized\n")
	fmt.Fprintf(writer, "Service account: %s\n\n", util.Bold("0x"+account.Address().String()))
	fmt.Fprintf(writer,
		"Start emulator by running: %s \nReset configuration using: %s.\n",
		util.Bold("flow emulator start"),
		util.Bold("flow init --reset"),
	)

	writer.Flush()
	return b.String()
}

// Oneliner show result as one liner grep friendly
func (r *InitResult) Oneliner() string {
	return ""
}
