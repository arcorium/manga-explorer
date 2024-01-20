package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetDefaultString(t *testing.T) {
	type args struct {
		real *string
		def  string
	}

	nonDefaultStr := "asdad"
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "Nil string",
			args: args{
				real: nil,
				def:  "test",
			},
		},
		{
			name: "Empty string",
			args: args{
				real: new(string),
				def:  "test",
			},
			expected: "test",
		},
		{
			name: "Non-Empty string",
			args: args{
				real: &nonDefaultStr,
				def:  "test",
			},
			expected: nonDefaultStr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDefaultString(tt.args.real, tt.args.def)
			if tt.args.real != nil {
				assert.Equal(t, *tt.args.real, tt.expected)
			}
		})
	}
}
