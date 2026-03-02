package main

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

type Test struct {
	name    string
	input   []string
	options UniqOptions
	expect  UniqResult
}

type UniqResult struct {
	output []string
	err    error
}

func TestUniq(t *testing.T) {
	tests := []Test{
		{
			name:    "Standard unique filtering (no flags)",
			input:   []string{"I love music.", "I love music.", "Thanks.", "I love music."},
			options: UniqOptions{},
			expect:  UniqResult{output: []string{"I love music.", "Thanks.", "I love music."}, err: nil},
		},
		{
			name:    "Count lines (-c)",
			input:   []string{"apple", "apple", "banana"},
			options: UniqOptions{CountLines: true},
			expect:  UniqResult{output: []string{"2 apple", "1 banana"}, err: nil},
		},
		{
			name:    "Duplicates only (-d)",
			input:   []string{"apple", "apple", "banana", "cherry", "cherry"},
			options: UniqOptions{DuplicatesOnly: true},
			expect:  UniqResult{output: []string{"apple", "cherry"}, err: nil},
		},
		{
			name:    "Unique only (-u)",
			input:   []string{"apple", "apple", "banana", "cherry"},
			options: UniqOptions{UniqueOnly: true},
			expect:  UniqResult{output: []string{"banana", "cherry"}, err: nil},
		},
		{
			name:    "Ignore case (-i)",
			input:   []string{"APPLE", "apple", "Banana"},
			options: UniqOptions{IgnoreCase: true},
			expect:  UniqResult{output: []string{"APPLE", "Banana"}, err: nil},
		},
		{
			name:    "Ignore fields (-f)",
			input:   []string{"row1 apple", "row2 apple", "row3 banana"},
			options: UniqOptions{IgnoreFields: 1},
			expect:  UniqResult{output: []string{"row1 apple", "row3 banana"}, err: nil},
		},
		{
			name:    "Ignore more fields than exist (-f)",
			input:   []string{"row1 apple", "row2 apple", "row3 banana"},
			options: UniqOptions{IgnoreFields: 4},
			expect:  UniqResult{output: []string{"row1 apple"}, err: nil},
		},
		{
			name:    "Ignore characters (-s)",
			input:   []string{"1abc", "2abc", "3def"},
			options: UniqOptions{IgnoreCharacters: 1},
			expect:  UniqResult{output: []string{"1abc", "3def"}, err: nil},
		},
		{
			name:    "Ignore more characters than exist (-s)",
			input:   []string{"1abc", "2abc", "3def"},
			options: UniqOptions{IgnoreCharacters: 9},
			expect:  UniqResult{output: []string{"1abc"}, err: nil},
		},
		{
			name:    "Invalid parameters (conflict -c and -d)",
			input:   []string{"apple", "apple"},
			options: UniqOptions{CountLines: true, DuplicatesOnly: true},
			expect:  UniqResult{output: []string{}, err: errors.New("Invalid options")},
		},
		{
			name:    "Empty input",
			input:   []string{},
			options: UniqOptions{},
			expect:  UniqResult{output: []string{}, err: nil},
		},
		{
			name:    "Ignore UTF-8 characters",
			input:   []string{"Привет", "Gривет", "食ривет"},
			options: UniqOptions{IgnoreCharacters: 1},
			expect:  UniqResult{output: []string{"Привет"}, err: nil},
		},
		{
			name:    "Empty strings",
			input:   []string{"", "", "", "123", "123", "", ""},
			options: UniqOptions{},
			expect:  UniqResult{output: []string{"", "123", ""}, err: nil},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result, err := uniq(test.input, test.options)
			actual := UniqResult{
				output: result,
				err:    err,
			}

			require.Equal(t, test.expect, actual)
		})
	}
}
