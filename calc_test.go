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

func TestUniq(t *testing.T) {
	tests := []Test{
		{
			name:        "Basic calculation",
			input:       "1+2",
			output:      3,
			expectError: false,
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
