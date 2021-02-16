/*
 * Flow CLI
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
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
	"testing"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/assert"
)

func generateComplexConfig() Config {
	return Config{
		Contracts: Contracts{{
			Name:    "NonFungibleToken",
			Source:  "../hungry-kitties/cadence/contracts/NonFungibleToken.cdc",
			Network: "emulator",
		}, {
			Name:    "FungibleToken",
			Source:  "../hungry-kitties/cadence/contracts/FungibleToken.cdc",
			Network: "emulator",
		}, {
			Name:    "Kibble",
			Source:  "./cadence/kibble/contracts/Kibble.cdc",
			Network: "emulator",
		}, {
			Name:    "KittyItems",
			Source:  "./cadence/kittyItems/contracts/KittyItems.cdc",
			Network: "emulator",
		}, {
			Name:    "KittyItemsMarket",
			Source:  "./cadence/kittyItemsMarket/contracts/KittyItemsMarket.cdc",
			Network: "emulator",
		}, {
			Name:    "KittyItemsMarket",
			Source:  "0x123123123",
			Network: "testnet",
		}},
		Deploys: Deploys{{
			Network:   "emulator",
			Account:   "emulator-account",
			Contracts: []string{"KittyItems", "KittyItemsMarket"},
		}, {
			Network:   "emulator",
			Account:   "account-4",
			Contracts: []string{"FungibleToken", "NonFungibleToken", "Kibble", "KittyItems", "KittyItemsMarket"},
		}, {
			Network:   "testnet",
			Account:   "account-2",
			Contracts: []string{"FungibleToken", "NonFungibleToken", "Kibble", "KittyItems"},
		}},
		Accounts: Accounts{{
			Name:    "emulator-account",
			Address: flow.ServiceAddress(flow.Emulator),
			ChainID: flow.Emulator,
			Keys: []AccountKey{{
				Type:     KeyTypeHex,
				Index:    0,
				SigAlgo:  crypto.ECDSA_P256,
				HashAlgo: crypto.SHA3_256,
				Context: map[string]string{
					"privateKey": "dd72967fd2bd75234ae9037dd4694c1f00baad63a10c35172bf65fbb8ad74b47",
				},
			}},
		}, {
			Name:    "account-2",
			Address: flow.HexToAddress("2c1162386b0a245f"),
			ChainID: flow.Emulator,
			Keys: []AccountKey{{
				Type:     KeyTypeHex,
				Index:    0,
				SigAlgo:  crypto.ECDSA_P256,
				HashAlgo: crypto.SHA3_256,
				Context: map[string]string{
					"privateKey": "dd72967fd2bd75234ae9037dd4694c1f00baad63a10c35172bf65fbb8ad74b47",
				},
			}},
		}, {
			Name:    "account-4",
			Address: flow.HexToAddress("f8d6e0586b0a20c1"),
			ChainID: flow.Emulator,
			Keys: []AccountKey{{
				Type:     KeyTypeHex,
				Index:    0,
				SigAlgo:  crypto.ECDSA_P256,
				HashAlgo: crypto.SHA3_256,
				Context: map[string]string{
					"privateKey": "dd72967fd2bd75234ae9037dd4694c1f00baad63a10c35172bf65fbb8ad74b47",
				},
			}},
		}},
		Networks: Networks{{
			Name:    "emulator",
			Host:    "127.0.0.1.3569",
			ChainID: flow.Emulator,
		}},
	}
}

func Test_GetContractsForNetworkComplex(t *testing.T) {
	conf := generateComplexConfig()
	kitty := conf.Contracts.GetByName("KittyItems")
	market := conf.Contracts.GetByName("KittyItemsMarket")

	assert.Equal(t, "KittyItems", kitty.Name)
	assert.Equal(t, "./cadence/kittyItemsMarket/contracts/KittyItemsMarket.cdc", market.Source)
}

func Test_GetContractsByNameAndNetworkComplex(t *testing.T) {
	conf := generateComplexConfig()
	market := conf.Contracts.GetByNameAndNetwork("KittyItemsMarket", "testnet")

	assert.Equal(t, "0x123123123", market.Source)
}

func Test_GetContractsByNetworkComplex(t *testing.T) {
	conf := generateComplexConfig()
	contracts := conf.Contracts.GetByNetwork("emulator")

	assert.Equal(t, 5, len(contracts))
	assert.Equal(t, "NonFungibleToken", contracts[0].Name)
	assert.Equal(t, "FungibleToken", contracts[1].Name)
	assert.Equal(t, "Kibble", contracts[2].Name)
	assert.Equal(t, "KittyItems", contracts[3].Name)
	assert.Equal(t, "KittyItemsMarket", contracts[4].Name)
}

func Test_GetAccountByNameComplex(t *testing.T) {
	conf := generateComplexConfig()
	acc := conf.Accounts.GetByName("account-4")

	fmt.Println(acc)

	assert.Equal(t, "f8d6e0586b0a20c1", acc.Address.String())
}

func Test_GetAccountByAddressComplex(t *testing.T) {
	conf := generateComplexConfig()
	acc1 := conf.Accounts.GetByAddress("0xf8d6e0586b0a20c1")
	acc2 := conf.Accounts.GetByAddress("2c1162386b0a245f")

	assert.Equal(t, "account-4", acc1.Name)
	assert.Equal(t, "account-2", acc2.Name)
}

func Test_GetDeploysByNetworkComplex(t *testing.T) {
	conf := generateComplexConfig()
	deploys := conf.Deploys.GetByAccountAndNetwork("account-2", "emulator")

	assert.Equal(t, deploys[0].Contracts, []string{"FungibleToken", "NonFungibleToken", "Kibble", "KittyItems"})
}

func Test_GetNetworkByNameComplex(t *testing.T) {
	conf := generateComplexConfig()
	network := conf.Networks.GetByName("emulator")

	assert.Equal(t, network.Host, "127.0.0.1.3569")
}