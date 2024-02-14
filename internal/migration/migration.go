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

package migration

import (
	"fmt"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flowkit"
	"github.com/onflow/flowkit/accounts"
)

// address of the migration contract on each network
var migrationContractStagingAddress = map[string]string{
	"testnet":   "0x2ceae959ed1a7e7a",
	"crescendo": "0x27b2302520211b67",
}

// MigrationContractStagingAddress returns the address of the migration contract on the given network
func MigrationContractStagingAddress(network string) flow.Address {
	return flow.HexToAddress(migrationContractStagingAddress[network])
}

func getAccountByContractName(state *flowkit.State, contractName string, network string) (*accounts.Account, error) {
	deployments := state.Deployments().ByNetwork(network)
	var accountName string
	for _, d := range deployments {
		for _, c := range d.Contracts {
			if c.Name == contractName {
				accountName = d.Account
				break
			}
		}
	}
	if accountName == "" {
		return nil, fmt.Errorf("contract not found in state")
	}

	accs := state.Accounts()
	if accs == nil {
		return nil, fmt.Errorf("no accounts found in state")
	}

	var account *accounts.Account
	for _, a := range *accs {
		if accountName == a.Name {
			account = &a
			break
		}
	}
	if account == nil {
		return nil, fmt.Errorf("account %s not found in state", accountName)
	}

	return account, nil
}