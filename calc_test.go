package main

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type Test struct {
	name     string
	input    string
	expected CalcResult
}

type CalcResult struct {
	output float64
	err    error
}

func TestCalc(t *testing.T) {
	tests := []Test{
		{
			name:     "Basic calculation",
			input:    "1+2",
			expected: CalcResult{output: 3, err: nil},
		},
		{
			name:     "Operator precedence",
			input:    "2 + 3 * 4",
			expected: CalcResult{output: 14, err: nil},
		},
		{
			name:     "Parentheses precedence",
			input:    "(2 + 3) * 4",
			expected: CalcResult{output: 20, err: nil},
		},
		{
			name:     "Decimals",
			input:    " 10.5 / 2 + 0.25 ",
			expected: CalcResult{output: 5.5, err: nil},
		},
		{
			name:     "Nesting",
			input:    "100 / (2 * (10 - 5))",
			expected: CalcResult{output: 10, err: nil},
		},
		{
			name:     "Stacked parentheses",
			input:    "((((3))) * (3))",
			expected: CalcResult{output: 9, err: nil},
		},
		{
			name:     "Empty parentheses",
			input:    "() + 3 - ()",
			expected: CalcResult{output: 0, err: errors.New("Empty parentheses")},
		},
		{
			name:     "Implicit multiplication",
			input:    "(2)(2)",
			expected: CalcResult{output: 0, err: errors.New("Too many values in the expression")},
		},
		{
			name:     "Division by zero",
			input:    "10 / 0",
			expected: CalcResult{output: 0, err: errors.New("Division by zero")},
		},
		{
			name:     "Mismatched parentheses (extra opening)",
			input:    "((2 + 2)",
			expected: CalcResult{output: 0, err: fmt.Errorf("Mismatched parentheses: unclosed '%c'", groupOpen)},
		},
		{
			name:     "Mismatched parentheses (extra closing)",
			input:    "(2 + 2))",
			expected: CalcResult{output: 0, err: fmt.Errorf("Mismatched parentheses: unexpected '%c'", groupClose)},
		},
		{
			name:     "Mismatched parentheses (begin with closing)",
			input:    "))((2) * 3)",
			expected: CalcResult{output: 0, err: fmt.Errorf("Mismatched parentheses: unexpected '%c'", groupClose)},
		},
		{
			name:     "Incomplete expression (trailing operator)",
			input:    "2 + 2 *",
			expected: CalcResult{output: 0, err: errors.New("Not enough values to use with '+'")},
		},
		{
			name:     "Incomplete expression (not enought values)",
			input:    "2 ** 2",
			expected: CalcResult{output: 0, err: errors.New("Not enough values to use with '*'")},
		},
		{
			name:     "Incomplete expression (too many values)",
			input:    "2 2",
			expected: CalcResult{output: 0, err: errors.New("Too many values in the expression")},
		},
		{
			name:     "Invalid character",
			input:    "2 ^ 3",
			expected: CalcResult{output: 0, err: errors.New("Unexpected character '^'")},
		},
		{
			name:     "Unary operation",
			input:    "-5 + 3",
			expected: CalcResult{output: 0, err: errors.New("Not enough values to use with '-'")},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result, err := evaluate(test.input)
			actual := CalcResult{
				output: result,
				err:    err,
			}

			require.Equal(t, test.expected, actual)
		})
	}
}
