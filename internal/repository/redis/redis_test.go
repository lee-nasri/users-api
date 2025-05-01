package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_REPO_escapeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no special characters",
			input:    "hello world",
			expected: "hello\\ world",
		},
		{
			name:     "string with special characters",
			input:    "hello@world!",
			expected: "hello\\@world\\!",
		},
		{
			name:     "string with backslashes",
			input:    "hello\\world",
			expected: "hello\\\\world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_REPO_sliceToInterface(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		wantResp []interface{}
		wantErr  error
	}{
		{
			name:     "slice of strings",
			input:    []string{"hello", "world"},
			wantResp: []interface{}{"hello", "world"},
			wantErr:  nil,
		},
		{
			name:     "slice of integers",
			input:    []int{1, 2, 3},
			wantResp: []interface{}{1, 2, 3},
			wantErr:  nil,
		},
		{
			name:     "empty slice",
			input:    []string{},
			wantResp: []interface{}{},
			wantErr:  nil,
		},
		{
			name:     "nil slice",
			input:    nil,
			wantResp: nil,
			wantErr:  ErrInputNotSlice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sliceToInterface(tt.input)

			// Assert response
			assert.Equal(t, tt.wantResp, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
