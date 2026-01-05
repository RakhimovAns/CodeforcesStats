package utils

import "testing"

func TestEcsapeHTML(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeHTML(tt.args.text); got != tt.want {
				t.Errorf("EscapeMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}
