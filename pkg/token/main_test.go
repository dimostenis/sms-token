package token

import "testing"

func Test_read_stdin_ok(t *testing.T) {
	// my custom stdin
	const my_stdin string = "test input"

	inject_stdin(my_stdin)

	// call the function that reads from stdin
	result := read_stdin()

	// check if the result is correct
	if result != my_stdin {
		t.Errorf("read_stdin() = %q, want %s", result, my_stdin)
	}
}

func Test_str_to_int(t *testing.T) {
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
			if got := str_to_int(tt.args.s); got != tt.want {
				t.Errorf("str_to_int() = %v, want %v", got, tt.want)
			}
		})
	}
}
