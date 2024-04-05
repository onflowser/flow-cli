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
	"bytes"
	"fmt"

	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/spf13/cobra"

	"github.com/onflow/flowkit"
	"github.com/onflow/flowkit/config"
	"github.com/onflow/flowkit/output"

	"github.com/onflow/flow-cli/internal/command"
	"github.com/onflow/flow-cli/internal/util"
)

type flagsInit struct {
	ServicePrivateKey  string `flag:"service-private-key" info:"Service account private key"`
	ServiceKeySigAlgo  string `default:"ECDSA_P256" flag:"service-sig-algo" info:"Service account key signature algorithm"`
	ServiceKeyHashAlgo string `default:"SHA3_256" flag:"service-hash-algo" info:"Service account key hash algorithm"`
	Reset              bool   `default:"false" flag:"reset" info:"Reset configuration file"`
	Global             bool   `default:"false" flag:"global" info:"Initialize global user configuration"`
}

var InitFlag = flagsInit{}

var initCommand = &command.Command{
	Cmd: &cobra.Command{
		Use:   "init",
		Short: "Initialize a new configuration",
	},
	Flags: &InitFlag,
	Run:   Initialise,
}

// InitConfigParameters holds all necessary parameters for initializing the configuration.
type InitConfigParameters struct {
	ServicePrivateKey  string
	ServiceKeySigAlgo  string
	ServiceKeyHashAlgo string
	Reset              bool
	Global             bool
}

// InitializeConfiguration creates the Flow configuration json file based on the provided parameters.
func InitializeConfiguration(params InitConfigParameters, logger output.Logger, readerWriter flowkit.ReaderWriter) (*flowkit.State, error) {
	logger.Info("⚠️Notice: for starting a new project prefer using 'flow setup'.")

	sigAlgo := crypto.StringToSignatureAlgorithm(params.ServiceKeySigAlgo)
	if sigAlgo == crypto.UnknownSignatureAlgorithm {
		return nil, fmt.Errorf("invalid signature algorithm: %s", params.ServiceKeySigAlgo)
	}

	hashAlgo := crypto.StringToHashAlgorithm(params.ServiceKeyHashAlgo)
	if hashAlgo == crypto.UnknownHashAlgorithm {
		return nil, fmt.Errorf("invalid hash algorithm: %s", params.ServiceKeyHashAlgo)
	}

	state, err := flowkit.Init(readerWriter, sigAlgo, hashAlgo)
	if err != nil {
		return nil, err
	}

	if params.ServicePrivateKey != "" {
		privateKey, err := crypto.DecodePrivateKeyHex(sigAlgo, params.ServicePrivateKey)
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %w", err)
		}

		state.SetEmulatorKey(privateKey)
	}

	path := config.DefaultPath
	if params.Global {
		path = config.GlobalPath()
	}

	if config.Exists(path) && !params.Reset {
		return nil, fmt.Errorf(
			"configuration already exists at: %s, if you want to reset configuration use the reset flag",
			path,
		)
	}

	err = state.Save(path)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func Initialise(
	_ []string,
	_ command.GlobalFlags,
	logger output.Logger,
	readerWriter flowkit.ReaderWriter,
	_ flowkit.Services,
) (command.Result, error) {
	params := InitConfigParameters{
		ServicePrivateKey:  InitFlag.ServicePrivateKey,
		ServiceKeySigAlgo:  InitFlag.ServiceKeySigAlgo,
		ServiceKeyHashAlgo: InitFlag.ServiceKeyHashAlgo,
		Reset:              InitFlag.Reset,
		Global:             InitFlag.Global,
	}
	state, err := InitializeConfiguration(params, logger, readerWriter)
	if err != nil {
		return nil, err
	}

	return &InitResult{State: state}, nil
}

type InitResult struct {
	*flowkit.State
}

func (r *InitResult) JSON() any {
	return r
}

func (r *InitResult) String() string {
	var b bytes.Buffer
	writer := util.CreateTabWriter(&b)
	account, _ := r.State.EmulatorServiceAccount()

	_, _ = fmt.Fprintf(writer, "Configuration initialized\n")
	_, _ = fmt.Fprintf(writer, "Service account: %s\n\n", output.Bold("0x"+account.Address.String()))
	_, _ = fmt.Fprintf(writer,
		"Start emulator by running: %s \nReset configuration using: %s\n",
		output.Bold("'flow emulator'"),
		output.Bold("'flow init --reset'"),
	)

	_ = writer.Flush()
	return b.String()
}

func (r *InitResult) Oneliner() string {
	return ""
}
