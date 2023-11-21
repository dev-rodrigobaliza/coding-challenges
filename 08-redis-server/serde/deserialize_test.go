package serde_test

import (
	"errors"
	"reflect"
	"rs/serde"
	"testing"
)

func TestDeserialize(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    serde.Command
		wantErr error
	}{
		{
			name: "empty string",
			args: args{
				str: "",
			},
			want:    serde.Command{},
			wantErr: serde.ErrInvalidString,
		},
		{
			name: "invalid command type",
			args: args{
				str: "&invalidcommandtype" + serde.End,
			},
			want:    serde.Command{},
			wantErr: serde.ErrInvalidCommandType,
		},
		{
			name: "invalid error 1",
			args: args{
				str: "-invalidcommanderror" + serde.End,
			},
			want:    serde.Command{},
			wantErr: serde.ErrInvalidCommandError,
		},
		{
			name: "invalid error 2",
			args: args{
				str: "-errinvalidcommanderror" + serde.End,
			},
			want:    serde.Command{},
			wantErr: serde.ErrInvalidCommandError,
		},
		{
			name: "invalid error 3",
			args: args{
				str: "-err" + serde.End,
			},
			want:    serde.Command{},
			wantErr: serde.ErrInvalidCommandError,
		},
		{
			name: "valid error 1",
			args: args{
				str: "-err error message" + serde.End,
			},
			want: serde.Command{
				Type:  serde.Error,
				Value: "",
				Array: nil,
				Error: &serde.CommandError{
					Type:    "ERR",
					Message: "error message",
				},
			},
			wantErr: nil,
		},
		{
			name: "valid error 2",
			args: args{
				str: "-ERR Error message" + serde.End,
			},
			want: serde.Command{
				Type:  serde.Error,
				Value: "",
				Array: nil,
				Error: &serde.CommandError{
					Type:    "ERR",
					Message: "Error message",
				},
			},
			wantErr: nil,
		},
		{
			name: "empty simple string",
			args: args{
				str: "+" + serde.End,
			},
			want: serde.Command{
				Type:  serde.SimpleString,
				Value: "",
			},
			wantErr: nil,
		},
		{
			name: "simple string",
			args: args{
				str: "+test" + serde.End,
			},
			want: serde.Command{
				Type:  serde.SimpleString,
				Value: "test",
			},
			wantErr: nil,
		},
		{
			name: "empty bulk string",
			args: args{
				str: "$0" + serde.End + serde.End,
			},
			want: serde.Command{
				Type:  serde.BulkString,
				Value: "",
			},
			wantErr: nil,
		},
		{
			name: "bulk string wrong size",
			args: args{
				str: "$0" + serde.End + "bulk string" + serde.End,
			},
			want:    serde.Command{},
			wantErr: serde.ErrInvalidString,
		},
		{
			name: "bulk string wrong size",
			args: args{
				str: "$1" + serde.End + serde.End,
			},
			want:    serde.Command{},
			wantErr: serde.ErrInvalidString,
		},
		{
			name: "bulk string",
			args: args{
				str: "$11" + serde.End + "bulk string" + serde.End,
			},
			want:    serde.Command{},
			wantErr: serde.ErrInvalidString,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serde.Deserialize(tt.args.str)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Deserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Deserialize() = %v, want %v", got, tt.want)
			}
		})
	}
}