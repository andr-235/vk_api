package transport

import (
	"testing"
)

func TestTransportError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *TransportError
		want string
	}{
		{
			name: "nil",
			err:  nil,
			want: "<nil>",
		},
		{
			name: "without message",
			err: &TransportError{
				Code: 123,
			},
			want: "vk api error 123",
		},
		{
			name: "with message",
			err: &TransportError{
				Code:    5,
				Message: "User authorization failed",
			},
			want: "vk api error 5: User authorization failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
