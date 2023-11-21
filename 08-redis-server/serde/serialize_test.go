package serde_test

import (
	"errors"
	"rs/serde"
	"testing"
)

func TestSerialize(t *testing.T) {
	type args struct {
		cmd serde.Command
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "invalid command",
			args: args{
				serde.Command{
					Type:  serde.Unknown,
					Value: "",
					Array: nil,
					Error: nil,
				},
			},
			want:    "",
			wantErr: serde.ErrInvalidCommandType,
		},
		{
			name: "invalid command error",
			args: args{
				serde.Command{
					Type:  serde.Error,
					Value: "",
					Array: nil,
					Error: nil,
				},
			},
			want:    "",
			wantErr: serde.ErrInvalidCommandError,
		},
		{
			name: "valid command error",
			args: args{
				serde.Command{
					Type:  serde.Error,
					Value: "",
					Array: nil,
					Error: &serde.CommandError{
						Type:    "TEST",
						Message: "some error for test",
					},
				},
			},
			want:    "-TEST some error for test" + serde.End,
			wantErr: nil,
		},
		{
			name: "command simple string empty",
			args: args{
				serde.Command{
					Type:  serde.SimpleString,
					Value: "",
					Array: nil,
					Error: nil,
				},
			},
			want:    "+" + serde.End,
			wantErr: nil,
		},
		{
			name: "command simple string",
			args: args{
				serde.Command{
					Type:  serde.SimpleString,
					Value: "simple string test",
					Array: nil,
					Error: nil,
				},
			},
			want:    "+simple string test" + serde.End,
			wantErr: nil,
		},
		{
			name: "command bulk string empty",
			args: args{
				serde.Command{
					Type:  serde.BulkString,
					Value: "",
					Array: nil,
					Error: nil,
				},
			},
			want:    "$0" + serde.End + serde.End,
			wantErr: nil,
		},
		{
			name: "command bulk string",
			args: args{
				serde.Command{
					Type:  serde.BulkString,
					Value: "bulk string test",
					Array: nil,
					Error: nil,
				},
			},
			want:    "$16" + serde.End + "bulk string test" + serde.End,
			wantErr: nil,
		},
		{
			name: "command integer empty",
			args: args{
				serde.Command{
					Type:  serde.Integer,
					Value: "",
					Array: nil,
					Error: nil,
				},
			},
			want:    "",
			wantErr: serde.ErrInvalidCommandValue,
		},
		{
			name: "command integer invalid",
			args: args{
				serde.Command{
					Type:  serde.Integer,
					Value: "a",
					Array: nil,
					Error: nil,
				},
			},
			want:    "",
			wantErr: serde.ErrInvalidCommandValue,
		},
		{
			name: "command integer",
			args: args{
				serde.Command{
					Type:  serde.Integer,
					Value: "1",
					Array: nil,
					Error: nil,
				},
			},
			want:    ":1" + serde.End,
			wantErr: nil,
		},
		{
			name: "command integer positive",
			args: args{
				serde.Command{
					Type:  serde.Integer,
					Value: "+1",
					Array: nil,
					Error: nil,
				},
			},
			want:    ":+1" + serde.End,
			wantErr: nil,
		},
		{
			name: "command integer negative",
			args: args{
				serde.Command{
					Type:  serde.Integer,
					Value: "-1",
					Array: nil,
					Error: nil,
				},
			},
			want:    ":-1" + serde.End,
			wantErr: nil,
		},
		{
			name: "command array empty",
			args: args{
				serde.Command{
					Type:  serde.Array,
					Value: "",
					Array: nil,
					Error: nil,
				},
			},
			want:    "*0" + serde.End,
			wantErr: nil,
		},
		{
			name: "command array one string",
			args: args{
				serde.Command{
					Type:  serde.Array,
					Value: "",
					Array: []serde.CommandArray{
						{
							Type:  serde.SimpleString,
							Value: "one",
							Array: nil,
						},
					},
					Error: nil,
				},
			},
			want:    "*1" + serde.End + "+one" + serde.End,
			wantErr: nil,
		},
		{
			name: "command array two strings",
			args: args{
				serde.Command{
					Type:  serde.Array,
					Value: "",
					Array: []serde.CommandArray{
						{
							Type:  serde.SimpleString,
							Value: "one",
							Array: nil,
						},
						{
							Type:  serde.SimpleString,
							Value: "two",
							Array: nil,
						},
					},
					Error: nil,
				},
			},
			want:    "*2" + serde.End + "+one" + serde.End + "+two" + serde.End,
			wantErr: nil,
		},
		{
			name: "command array one integer",
			args: args{
				serde.Command{
					Type:  serde.Array,
					Value: "",
					Array: []serde.CommandArray{
						{
							Type:  serde.Integer,
							Value: "1",
							Array: nil,
						},
					},
					Error: nil,
				},
			},
			want:    "*1" + serde.End + ":1" + serde.End,
			wantErr: nil,
		},
		{
			name: "command array two integers",
			args: args{
				serde.Command{
					Type:  serde.Array,
					Value: "",
					Array: []serde.CommandArray{
						{
							Type:  serde.Integer,
							Value: "+1",
							Array: nil,
						},
						{
							Type:  serde.Integer,
							Value: "-1",
							Array: nil,
						},
					},
					Error: nil,
				},
			},
			want:    "*2" + serde.End + ":+1" + serde.End + ":-1" + serde.End,
			wantErr: nil,
		},
		{
			name: "command array mixed (string and integer)",
			args: args{
				serde.Command{
					Type:  serde.Array,
					Value: "",
					Array: []serde.CommandArray{
						{
							Type:  serde.SimpleString,
							Value: "string",
							Array: nil,
						},
						{
							Type:  serde.Integer,
							Value: "1",
							Array: nil,
						},
					},
					Error: nil,
				},
			},
			want:    "*2" + serde.End + "+string" + serde.End + ":1" + serde.End,
			wantErr: nil,
		},
		{
			name: "command array nested",
			args: args{
				serde.Command{
					Type:  serde.Array,
					Value: "",
					Array: []serde.CommandArray{
						{
							Type:  serde.Array,
							Value: "",
							Array: []serde.CommandArray{
								{
									Type:  serde.SimpleString,
									Value: "nested",
									Array: nil,
								},
							},
						},
					},
					Error: nil,
				},
			},
			want:    "*1" + serde.End + "*1" + serde.End + "+nested" + serde.End,
			wantErr: nil,
		},
		{
			name: "command array nested and mixed",
			args: args{
				serde.Command{
					Type:  serde.Array,
					Value: "",
					Array: []serde.CommandArray{
						{
							Type:  serde.Array,
							Value: "",
							Array: []serde.CommandArray{
								{
									Type:  serde.SimpleString,
									Value: "nested",
									Array: nil,
								},
							},
						},
						{
							Type:  serde.Integer,
							Value: "1",
							Array: nil,
						},
					},
					Error: nil,
				},
			},
			want:    "*2" + serde.End + "*1" + serde.End + "+nested" + serde.End + ":1" + serde.End,
			wantErr: nil,
		},
		{
			name: "command array nested, mixed and invalid",
			args: args{
				serde.Command{
					Type:  serde.Array,
					Value: "",
					Array: []serde.CommandArray{
						{
							Type:  serde.Array,
							Value: "",
							Array: []serde.CommandArray{
								{
									Type:  serde.SimpleString,
									Value: "nested",
									Array: nil,
								},
							},
						},
						{
							Type:  serde.Integer,
							Value: "",
							Array: nil,
						},
					},
					Error: nil,
				},
			},
			want:    "",
			wantErr: serde.ErrInvalidCommandValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serde.Serialize(tt.args.cmd)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}
