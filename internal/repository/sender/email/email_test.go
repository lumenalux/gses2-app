package email

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"gses2-app/internal/core/port"
	"gses2-app/internal/repository/sender/smtp"
)

var (
	errDialerError  = errors.New("dialer error")
	errFactoryError = errors.New("factory error")
)

func TestSendExchangeRate(t *testing.T) {
	tests := []struct {
		name         string
		emails       []string
		exchangeRate port.Rate
		dialer       smtp.TLSConnectionDialer
		factory      smtp.SMTPClientFactory
		expectedErr  error
	}{
		{
			name:         "Successful SendExchangeRate",
			emails:       []string{"test@example.com"},
			exchangeRate: 10.5,
			dialer:       &smtp.StubDialer{},
			factory:      &smtp.StubSMTPClientFactory{Client: &smtp.StubSMTPClient{}},
			expectedErr:  nil,
		},
		{
			name:         "Failed due to dialer error",
			emails:       []string{"test@example.com"},
			exchangeRate: 10.5,
			dialer:       &smtp.StubDialer{Err: errDialerError},
			factory:      &smtp.StubSMTPClientFactory{Client: &smtp.StubSMTPClient{}},
			expectedErr:  errDialerError,
		},
		{
			name:         "Failed due to factory error",
			emails:       []string{"test@example.com"},
			exchangeRate: 10.5,
			dialer:       &smtp.StubDialer{},
			factory: &smtp.StubSMTPClientFactory{
				Client: &smtp.StubSMTPClient{},
				Err:    errFactoryError,
			},
			expectedErr: errFactoryError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := &EmailSenderConfig{}
			service, err := NewProvider(config, tt.dialer, tt.factory)

			require.Equal(t, tt.expectedErr, err)

			if tt.expectedErr != nil {
				return
			}

			users := convertEmailsToUsers(tt.emails)
			err = service.SendExchangeRate(tt.exchangeRate, users)

			require.NoError(t, err, "SendExchangeRate() unexpected error = %v", err)
		})
	}
}

func convertEmailsToUsers(emails []string) []port.User {
	users := make([]port.User, len(emails))

	for i, email := range emails {
		users[i] = port.User{Email: email}
	}

	return users
}
