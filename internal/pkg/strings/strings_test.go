package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUnique(t *testing.T) {
	t.Parallel()

	type args struct {
		length uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name:    "success",
			args:    args{length: 10},
			wantErr: false,
			err:     nil,
		},
		{
			name:    "0 length",
			args:    args{length: 0},
			wantErr: true,
			err:     errZeroLength,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUnique(tt.args.length)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)

			assert.NotEqual(t, "", got)
			assert.Equal(t, tt.args.length, uint(len(got)))
		})
	}
}

func TestNewRandom(t *testing.T) {
	t.Parallel()

	type args struct {
		length uint
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{length: 35},
		},
		{
			name: "0 length",
			args: args{length: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRandom(tt.args.length)
			assert.Equal(t, tt.args.length, uint(len(got)))
		})
	}
}

func TestNewSpecialRandom(t *testing.T) {
	t.Parallel()

	type args struct {
		length uint
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{length: 33},
		},
		{
			name: "0 length",
			args: args{length: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSpecialRandom(tt.args.length)
			assert.Equal(t, tt.args.length, uint(len(got)))
		})
	}
}
