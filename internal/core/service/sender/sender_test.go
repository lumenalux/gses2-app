package sender

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"gses2-app/internal/core/port"
)

type StubProvider struct {
	Err error
}

func (tp *StubProvider) SendExchangeRate(
	rate port.Rate,
	subscribers []port.User,
) error {
	return tp.Err
}

var (
	errProvider = errors.New("provider error")
)

func TestSendExchangeRate(t *testing.T) {
	tests := []struct {
		name        string
		providerErr error
		expectedErr error
	}{
		{
			name:        "No error from provider",
			providerErr: nil,
			expectedErr: nil,
		},
		{
			name:        "Error from provider",
			providerErr: errProvider,
			expectedErr: errProvider,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			provider := &StubProvider{Err: tt.providerErr}
			service := NewService(provider)

			err := service.SendExchangeRate(1.23, port.User{Email: "subscriber"})

			require.Equal(t, tt.expectedErr, err)
		})
	}
}
