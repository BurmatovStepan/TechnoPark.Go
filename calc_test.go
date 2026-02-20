package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type Test struct {
	name        string
	input       string
	output      float64
	expectError bool
}

func TestCalc(t *testing.T) {
	tests := []Test{
		{
			name:        "Basic calculation",
			input:       "1+2",
			output:      3,
			expectError: false,
		},
		{
			name:        "Operator precedence",
			input:       "2 + 3 * 4",
			output:      14,
			expectError: false,
		},
		{
			name:        "Parentheses precedence",
			input:       "(2 + 3) * 4",
			output:      20,
			expectError: false,
		},
		{
			name:        "Decimals",
			input:       " 10.5 / 2 + 0.25 ",
			output:      5.5,
			expectError: false,
		},
		{
			name:        "Nesting",
			input:       "100 / (2 * (10 - 5))",
			output:      10,
			expectError: false,
		},
		{
			name:        "Stacked parentheses",
			input:       "((((3))) * (3))",
			output:      9,
			expectError: false,
		},
		{
			name:        "Empty parentheses",
			input:       "() + 3 - ()",
			expectError: true,
		},
		{
			name:        "Implicit multiplication",
			input:       "(2)(2)",
			expectError: true,
		},
		{
			name:        "Division by zero",
			input:       "10 / 0",
			expectError: true,
		},
		{
			name:        "Mismatched parentheses (extra opening)",
			input:       "((2 + 2)",
			expectError: true,
		},
		{
			name:        "Mismatched parentheses (extra closing)",
			input:       "(2 + 2))",
			expectError: true,
		},
		{
			name:        "Mismatched parentheses (begin with closing)",
			input:       "))((2) * 3)",
			expectError: true,
		},
		{
			name:        "Incomplete expression (trailing operator)",
			input:       "2 + 2 *",
			expectError: true,
		},
		{
			name:        "Incomplete expression (not enought values)",
			input:       "2 ** 2",
			expectError: true,
		},
		{
			name:        "Incomplete expression (too many values)",
			input:       "2 2",
			expectError: true,
		},
		{
			name:        "Invalid character",
			input:       "2 ^ 3",
			expectError: true,
		},
		{
			name:        "Unary operation",
			input:       "-5 + 3",
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := evaluate(test.input)

			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.output, result)
			}
		})
	}
}
