package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFileFormat(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "Empty string",
			args: args{
				filename: "",
			},
			want:    "",
			wantErr: NoFormatErr,
		},
		{
			name: "Normal",
			args: args{
				filename: "format.png",
			},
			want:    "png",
			wantErr: nil,
		},
		{
			name: "Non Format",
			args: args{
				filename: "format",
			},
			want:    "",
			wantErr: NoFormatErr,
		},
		{
			name: "Multiple dot",
			args: args{
				filename: "format1.format2.png",
			},
			want:    "png",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFileFormat(tt.args.filename)
			if tt.wantErr != nil {
				assert.Equalf(t, err, tt.wantErr, "GetFileFormat(%v)", tt.args.filename)
				return
			}
			assert.Equalf(t, tt.want, got, "GetFileFormat(%v)", tt.args.filename)
		})
	}
}
