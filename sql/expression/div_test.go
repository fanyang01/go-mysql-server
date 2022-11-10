// Copyright 2022 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package expression

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dolthub/go-mysql-server/sql"
)

func TestDiv(t *testing.T) {
	var floatTestCases = []struct {
		name        string
		left, right float64
		expected    string
		null        bool
	}{
		{"1 / 1", 1, 1, "1.0000", false},
		{"1 / 2", 1, 2, "0.5000", false},
		{"-1 / 1.0", -1, 1, "-1.0000", false},
		{"0 / 1234567890", 0, 12345677890, "0.0000", false},
		{"3.14159 / 3.0", 3.14159, 3.0, "1.047196667", false},
		{"1/0", 1, 0, "", true},
		{"-1/0", -1, 0, "", true},
		{"0/0", 0, 0, "", true},
	}

	for _, tt := range floatTestCases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewDiv(
				// The numbers are interpreted as Float64 without going through parser, so we lose precision here for 1.0
				NewLiteral(tt.left, sql.Float64),
				NewLiteral(tt.right, sql.Float64),
			).Eval(sql.NewEmptyContext(), sql.NewRow())
			require.NoError(t, err)
			if tt.null {
				assert.Equal(t, nil, result)
			} else {
				r, ok := result.(decimal.Decimal)
				assert.True(t, ok)
				assert.Equal(t, tt.expected, r.StringFixed(r.Exponent()*-1))
			}
		})
	}

	var intTestCases = []struct {
		name        string
		left, right int64
		expected    string
		null        bool
	}{
		{"1 / 1", 1, 1, "1.0000", false},
		{"-1 / 1", -1, 1, "-1.0000", false},
		{"0 / 1234567890", 0, 12345677890, "0.0000", false},
		{"1/0", 1, 0, "", true},
		{"0/0", 1, 0, "", true},
	}
	for _, tt := range intTestCases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewDiv(
				NewLiteral(tt.left, sql.Int64),
				NewLiteral(tt.right, sql.Int64),
			).Eval(sql.NewEmptyContext(), sql.NewRow())
			require.NoError(t, err)
			if tt.null {
				assert.Equal(t, nil, result)
			} else {
				r, ok := result.(decimal.Decimal)
				assert.True(t, ok)
				assert.Equal(t, tt.expected, r.StringFixed(r.Exponent()*-1))
			}
		})
	}

	var uintTestCases = []struct {
		name        string
		left, right uint64
		expected    string
		null        bool
	}{
		{"1 / 1", 1, 1, "1.0000", false},
		{"0 / 1234567890", 0, 12345677890, "0.0000", false},
		{"1/0", 1, 0, "", true},
		{"0/0", 1, 0, "", true},
	}
	for _, tt := range uintTestCases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewDiv(
				NewLiteral(tt.left, sql.Uint64),
				NewLiteral(tt.right, sql.Uint64),
			).Eval(sql.NewEmptyContext(), sql.NewRow())
			require.NoError(t, err)
			if tt.null {
				assert.Equal(t, nil, result)
			} else {
				r, ok := result.(decimal.Decimal)
				assert.True(t, ok)
				assert.Equal(t, tt.expected, r.StringFixed(r.Exponent()*-1))
			}
		})
	}
}