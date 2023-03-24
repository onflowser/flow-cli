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

package gateway

import (
	"context"
	"fmt"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-emulator/convert/sdk"
	"github.com/onflow/flow-emulator/server/backend"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	flowGo "github.com/onflow/flow-go/model/flow"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

type EmulatorKey struct {
	PublicKey crypto.PublicKey
	SigAlgo   crypto.SignatureAlgorithm
	HashAlgo  crypto.HashAlgorithm
}

type EmulatorGateway struct {
	emulator        *emulator.Blockchain
	backend         *backend.Backend
	ctx             context.Context
	logger          *logrus.Logger
	emulatorOptions []emulator.Option
}

func UnwrapStatusError(err error) error {
	return fmt.Errorf(status.Convert(err).Message())
}

func NewEmulatorGateway(key *EmulatorKey) *EmulatorGateway {
	return NewEmulatorGatewayWithOpts(key)
}

func NewEmulatorGatewayWithOpts(key *EmulatorKey, opts ...func(*EmulatorGateway)) *EmulatorGateway {

	gateway := &EmulatorGateway{
		ctx:             context.Background(),
		logger:          logrus.New(),
		emulatorOptions: []emulator.Option{},
	}
	for _, opt := range opts {
		opt(gateway)
	}

	gateway.emulator = newEmulator(key, gateway.emulatorOptions...)
	gateway.backend = backend.New(&zerolog.Logger{}, gateway.emulator)
	gateway.backend.EnableAutoMine()

	return gateway
}

func WithEmulatorOptions(options ...emulator.Option) func(g *EmulatorGateway) {
	return func(g *EmulatorGateway) {
		g.emulatorOptions = append(g.emulatorOptions, options...)
	}
}

func (g *EmulatorGateway) SetContext(ctx context.Context) {
	g.ctx = ctx
}

func newEmulator(key *EmulatorKey, emulatorOptions ...emulator.Option) *emulator.Blockchain {
	var opts []emulator.Option

	if key != nil {
		opts = append(opts, emulator.WithServicePublicKey(key.PublicKey, key.SigAlgo, key.HashAlgo))
	}

	opts = append(opts, emulatorOptions...)

	b, err := emulator.NewBlockchain(opts...)
	if err != nil {
		panic(err)
	}

	return b
}

func (g *EmulatorGateway) GetAccount(address flow.Address) (*flow.Account, error) {
	account, err := g.backend.GetAccount(g.ctx, address)
	if err != nil {
		return nil, UnwrapStatusError(err)
	}
	return account, nil
}

func (g *EmulatorGateway) SendSignedTransaction(tx *flow.Transaction) (*flow.Transaction, error) {
	err := g.backend.SendTransaction(context.Background(), *tx)
	if err != nil {
		return nil, UnwrapStatusError(err)
	}
	return tx, nil
}

func (g *EmulatorGateway) GetTransactionResult(ID flow.Identifier, waitSeal bool) (*flow.TransactionResult, error) {
	result, err := g.backend.GetTransactionResult(g.ctx, ID)
	if err != nil {
		return nil, UnwrapStatusError(err)
	}
	return result, nil
}

func (g *EmulatorGateway) GetTransaction(id flow.Identifier) (*flow.Transaction, error) {
	transaction, err := g.backend.GetTransaction(g.ctx, id)
	if err != nil {
		return nil, UnwrapStatusError(err)
	}
	return transaction, nil
}

func (g *EmulatorGateway) GetTransactionResultsByBlockID(blockID flow.Identifier) ([]*flow.TransactionResult, error) {
	// TODO: implement
	panic("GetTransactionResultsByBlockID not implemented")
}

func (g *EmulatorGateway) GetTransactionsByBlockID(blockID flow.Identifier) ([]*flow.Transaction, error) {
	// TODO: implement
	panic("GetTransactionResultsByBlockID not implemented")
}

func (g *EmulatorGateway) Ping() error {
	err := g.backend.Ping(g.ctx)
	if err != nil {
		return UnwrapStatusError(err)
	}
	return nil
}

func (g *EmulatorGateway) ExecuteScript(script []byte, arguments []cadence.Value) (cadence.Value, error) {

	args, err := cadenceValuesToMessages(arguments)
	if err != nil {
		return nil, UnwrapStatusError(err)
	}

	result, err := g.backend.ExecuteScriptAtLatestBlock(g.ctx, script, args)
	if err != nil {
		return nil, UnwrapStatusError(err)
	}

	value, err := messageToCadenceValue(result)
	if err != nil {
		return nil, UnwrapStatusError(err)
	}

	return value, nil
}

func (g *EmulatorGateway) GetLatestBlock() (*flow.Block, error) {
	block, _, err := g.backend.GetLatestBlock(g.ctx, true)
	if err != nil {
		return nil, UnwrapStatusError(err)
	}

	return convertBlock(block), nil
}

func cadenceValuesToMessages(values []cadence.Value) ([][]byte, error) {
	msgs := make([][]byte, len(values))
	for i, val := range values {
		msg, err := jsoncdc.Encode(val)
		if err != nil {
			return nil, fmt.Errorf("convert: %w", err)
		}
		msgs[i] = msg
	}
	return msgs, nil
}

func messageToCadenceValue(m []byte) (cadence.Value, error) {
	v, err := jsoncdc.Decode(nil, m)
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}

	return v, nil
}

func convertBlock(block *flowGo.Block) *flow.Block {
	return &flow.Block{
		BlockHeader: flow.BlockHeader{
			ID:        flow.Identifier(block.Header.ID()),
			ParentID:  flow.Identifier(block.Header.ParentID),
			Height:    block.Header.Height,
			Timestamp: block.Header.Timestamp,
		},
		BlockPayload: flow.BlockPayload{
			CollectionGuarantees: nil,
			Seals:                nil,
		},
	}
}

func (g *EmulatorGateway) GetEvents(
	eventType string,
	startHeight uint64,
	endHeight uint64,
) ([]flow.BlockEvents, error) {
	events := make([]flow.BlockEvents, 0)

	for height := startHeight; height <= endHeight; height++ {
		events = append(events, g.getBlockEvent(height, eventType))
	}

	return events, nil
}

func (g *EmulatorGateway) getBlockEvent(height uint64, eventType string) flow.BlockEvents {
	block, _, _ := g.backend.GetBlockByHeight(g.ctx, height)
	events, _ := g.backend.GetEventsForBlockIDs(g.ctx, eventType, []flow.Identifier{flow.Identifier(block.ID())})

	result := flow.BlockEvents{
		BlockID:        flow.Identifier(block.ID()),
		Height:         block.Header.Height,
		BlockTimestamp: block.Header.Timestamp,
		Events:         []flow.Event{},
	}

	for _, e := range events {
		if e.BlockID == block.ID() {
			result.Events, _ = sdk.FlowEventsToSDK(e.Events)
			return result
		}
	}

	return result
}

func (g *EmulatorGateway) GetCollection(id flow.Identifier) (*flow.Collection, error) {
	collection, err := g.backend.GetCollectionByID(g.ctx, id)
	if err != nil {
		return nil, UnwrapStatusError(err)
	}
	return collection, nil
}

func (g *EmulatorGateway) GetBlockByID(id flow.Identifier) (*flow.Block, error) {
	block, _, err := g.backend.GetBlockByID(g.ctx, id)
	if err != nil {
		return nil, UnwrapStatusError(err)
	}
	return convertBlock(block), nil
}

func (g *EmulatorGateway) GetBlockByHeight(height uint64) (*flow.Block, error) {
	block, _, err := g.backend.GetBlockByHeight(g.ctx, height)
	if err != nil {
		return nil, UnwrapStatusError(err)
	}
	return convertBlock(block), nil
}

func (g *EmulatorGateway) GetLatestProtocolStateSnapshot() ([]byte, error) {
	snapshot, err := g.backend.GetLatestProtocolStateSnapshot(g.ctx)
	if err != nil {
		return nil, UnwrapStatusError(err)
	}
	return snapshot, nil
}

// SecureConnection placeholder func to complete gateway interface implementation
func (g *EmulatorGateway) SecureConnection() bool {
	return false
}
