/*
 * Flow CLI
 *
 * Copyright 2022 Dapper Labs, Inc.
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

package test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/onflow/cadence/runtime/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-cli/internal/util"
	"github.com/onflow/flow-cli/pkg/flowkit/config"
	"github.com/onflow/flow-cli/pkg/flowkit/tests"
)

func TestExecutingTests(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		t.Parallel()
		_, state, _ := util.TestMocks(t)

		script := tests.TestScriptSimple
		testFiles := make(map[string][]byte, 0)
		testFiles[script.Filename] = script.Source
		results, _, err := testCode(testFiles, state, false)

		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.NoError(t, results[script.Filename][0].Error)
	})

	t.Run("simple failing", func(t *testing.T) {
		t.Parallel()
		_, state, _ := util.TestMocks(t)

		script := tests.TestScriptSimpleFailing
		testFiles := make(map[string][]byte, 0)
		testFiles[script.Filename] = script.Source
		results, _, err := testCode(testFiles, state, false)

		require.NoError(t, err)
		require.Len(t, results, 1)

		err = results[script.Filename][0].Error
		require.Error(t, err)
		assert.ErrorAs(t, err, &stdlib.AssertionError{})
	})

	t.Run("with import", func(t *testing.T) {
		t.Parallel()
		_, state, _ := util.TestMocks(t)

		c := config.Contract{
			Name:     tests.ContractHelloString.Name,
			Location: tests.ContractHelloString.Filename,
		}
		state.Contracts().AddOrUpdate(c)

		// Execute script
		script := tests.TestScriptWithImport
		testFiles := make(map[string][]byte, 0)
		testFiles[script.Filename] = script.Source
		results, _, err := testCode(testFiles, state, false)

		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.NoError(t, results[script.Filename][0].Error)
	})

	t.Run("with file read", func(t *testing.T) {
		t.Parallel()
		_, state, rw := util.TestMocks(t)

		_ = rw.WriteFile(
			tests.SomeFile.Filename,
			tests.SomeFile.Source,
			os.ModeTemporary,
		)

		// Execute script
		script := tests.TestScriptWithFileRead
		testFiles := make(map[string][]byte, 0)
		testFiles[script.Filename] = script.Source
		results, _, err := testCode(testFiles, state, false)

		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.NoError(t, results[script.Filename][0].Error)
	})

	t.Run("with code coverage", func(t *testing.T) {
		t.Parallel()

		// Setup
		_, state, _ := util.TestMocks(t)

		state.Contracts().AddOrUpdate(config.Contract{
			Name:     tests.ContractFooCoverage.Name,
			Location: tests.ContractFooCoverage.Filename,
		})

		// Execute script
		script := tests.TestScriptWithCoverage
		testFiles := make(map[string][]byte, 0)
		testFiles[script.Filename] = script.Source
		results, coverageReport, err := testCode(testFiles, state, true)

		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.NoError(t, results[script.Filename][0].Error)

		actual, err := json.Marshal(coverageReport)
		require.NoError(t, err)

		expected := `
		  {
		    "coverage": {
		      "S.FooContract.cdc": {
		        "line_hits": {
		          "14": 1,
		          "18": 10,
		          "19": 1,
		          "20": 9,
		          "21": 1,
		          "22": 8,
		          "23": 1,
		          "24": 7,
		          "25": 1,
		          "26": 6,
		          "27": 1,
		          "30": 5,
		          "31": 4,
		          "34": 1,
		          "6": 1
		        },
		        "missed_lines": [],
		        "statements": 15,
		        "percentage": "100.0%"
		      }
		    }
		  }
		`

		require.JSONEq(t, expected, string(actual))

		assert.Equal(
			t,
			"Coverage: 100.0% of statements",
			coverageReport.CoveredStatementsPercentage(),
		)
	})
}
