package sms

import "testing"

func TestStrToInt(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "normal",
			args: args{s: "123"},
			want: 123,
		},
		{
			name: "zero",
			args: args{s: "0"},
			want: 0,
		},
		{
			name: "negative",
			args: args{s: "-5"},
			want: -5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strToInt(tt.args.s); got != tt.want {
				t.Errorf("str_to_int() = %v, want %v", got, tt.want)
			}
		})
	}
}
