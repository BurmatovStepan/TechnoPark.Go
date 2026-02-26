package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type Test struct {
	name        string
	input       []string
	output      []string
	options     UniqOptions
	expectError bool
}

func TestUniq(t *testing.T) {
	tests := []Test{
		{
			name:        "Standard unique filtering (no flags)",
			input:       []string{"I love music.", "I love music.", "Thanks.", "I love music."},
			output:      []string{"I love music.", "Thanks.", "I love music."},
			options:     UniqOptions{},
			expectError: false,
		},
		{
			name:        "Count lines (-c)",
			input:       []string{"apple", "apple", "banana"},
			output:      []string{"2 apple", "1 banana"},
			options:     UniqOptions{CountLines: true},
			expectError: false,
		},
		{
			name:        "Duplicates only (-d)",
			input:       []string{"apple", "apple", "banana", "cherry", "cherry"},
			output:      []string{"apple", "cherry"},
			options:     UniqOptions{DuplicatesOnly: true},
			expectError: false,
		},
		{
			name:        "Unique only (-u)",
			input:       []string{"apple", "apple", "banana", "cherry"},
			output:      []string{"banana", "cherry"},
			options:     UniqOptions{UniqueOnly: true},
			expectError: false,
		},
		{
			name:        "Ignore case (-i)",
			input:       []string{"APPLE", "apple", "Banana"},
			output:      []string{"APPLE", "Banana"},
			options:     UniqOptions{IgnoreCase: true},
			expectError: false,
		},
		{
			name:        "Ignore fields (-f)",
			input:       []string{"row1 apple", "row2 apple", "row3 banana"},
			output:      []string{"row1 apple", "row3 banana"},
			options:     UniqOptions{IgnoreFields: 1},
			expectError: false,
		},
		{
			name:        "Ignore more fields than exist (-f)",
			input:       []string{"row1 apple", "row2 apple", "row3 banana"},
			output:      []string{"row1 apple"},
			options:     UniqOptions{IgnoreFields: 4},
			expectError: false,
		},
		{
			name:        "Ignore characters (-s)",
			input:       []string{"1abc", "2abc", "3def"},
			output:      []string{"1abc", "3def"},
			options:     UniqOptions{IgnoreCharacters: 1},
			expectError: false,
		},
		{
			name:        "Ignore more characters than exist (-s)",
			input:       []string{"1abc", "2abc", "3def"},
			output:      []string{"1abc"},
			options:     UniqOptions{IgnoreCharacters: 9},
			expectError: false,
		},
		{
			name:        "Invalid parameters (conflict -c and -d)",
			input:       []string{"apple", "apple"},
			output:      []string{},
			options:     UniqOptions{CountLines: true, DuplicatesOnly: true},
			expectError: true,
		},
		{
			name:        "Empty input",
			input:       []string{},
			output:      []string{},
			options:     UniqOptions{},
			expectError: false,
		},
		{
			name:        "Ignore UTF-8 characters",
			input:       []string{"Привет", "Gривет", "食ривет"},
			output:      []string{"Привет"},
			options:     UniqOptions{IgnoreCharacters: 1},
			expectError: false,
		},
		{
			name:        "Empty strings",
			input:       []string{"", "", "", "123", "123", "", ""},
			output:      []string{"", "123", ""},
			options:     UniqOptions{},
			expectError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := uniq(test.input, test.options)

			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.output, result)
			}
		})
	}
}
