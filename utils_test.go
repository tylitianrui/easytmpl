package easytmpl

import (
	"testing"
)

func TestIsBlank(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		b    []byte
		want bool
	}{
		{
			name: "",
			b:    []byte{},
			want: true,
		},
		{
			name: "",
			b:    []byte{'1'}[0:0],
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBlank(tt.b)
			if got != tt.want {
				t.Errorf("IsBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}
