package migration

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/onflow/cadence"
	flowsdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flowkit"
	"github.com/onflow/flowkit/output"
	"github.com/spf13/cobra"

	"github.com/onflow/flow-cli/internal/command"
	"github.com/onflow/flow-cli/internal/scripts"
)

var getStagedCodeflags = scripts.Flags{}

var getStagedCodeCommand = &command.Command{
	Cmd: &cobra.Command{
		Use:     "flow staged-code <CONTRACT_ADDRESS>",
		Short:   "returns back the staged code for a contract",
		Example: `flow staged-code 0xhello`,
		Args:    cobra.MinimumNArgs(1),
	},
	Flags: &getStagedCodeflags,
	RunS:  getStagedCode,
}

func getStagedCode(
	args []string,
	globalFlags command.GlobalFlags,
	_ output.Logger,
	flow flowkit.Services,
	state *flowkit.State,
) (command.Result, error) {
	scTempl, err := template.ParseFiles("./cadence/scripts/get_staged_code.cdc")
	if err != nil {
		return nil, fmt.Errorf("error loading staging contract file: %w", err)
	}

	// render transaction template with network
	var txScriptBuf bytes.Buffer
	if err := scTempl.Execute(
		&txScriptBuf,
		map[string]string{
			"MigrationContractStaging": MigrationContractStagingAddress[globalFlags.Network],
		}); err != nil {
		return nil, fmt.Errorf("error rendering staging contract template: %w", err)
	}

	contractAddress := args[0]

	caddr := cadence.NewAddress(flowsdk.HexToAddress(contractAddress))

	query := flowkit.ScriptQuery{}
	if getStagedCodeflags.BlockHeight != 0 {
		query.Height = getStagedCodeflags.BlockHeight
	} else if getStagedCodeflags.BlockID != "" {
		query.ID = flowsdk.HexToID(getStagedCodeflags.BlockID)
	} else {
		query.Latest = true
	}

	value, err := flow.ExecuteScript(
		context.Background(),
		flowkit.Script{
			Code: txScriptBuf.Bytes(),
			Args: []cadence.Value{caddr},
		},
		query,
	)
	if err != nil {
		return nil, err
	}

	return scripts.NewScriptResult(value), nil
}