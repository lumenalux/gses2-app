// Subscription integration contains integration tests for the subscription layer,
// specifically testing the integration between the subscription service and
// the user repository with storage

package integration

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"gses2-app/internal/core/port"
	"gses2-app/internal/core/service/subscription"
	"gses2-app/internal/repository/storage"
)

type SubscriptionTest struct {
	Name           string
	Subscribers    []port.User
	Action         func(service *subscription.Service, emails []port.User) error
	ExpectedError  error
	ExpectedResult []port.User
}

func TestSubscriptionServiceIntegration(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	storageCSV := storage.NewCSVStorage(tmpFile.Name())
	userRepository := port.NewUserRepository(storageCSV)
	service := subscription.NewService(userRepository)

	tests := []SubscriptionTest{
		{
			Name:        "Subscribe a new email",
			Subscribers: []port.User{{Email: "test1@example.com"}},
			Action: func(service *subscription.Service, subscribers []port.User) error {
				return service.Subscribe(&subscribers[0])
			},
		},
		{
			Name:        "Subscribe an already subscribed email",
			Subscribers: []port.User{{Email: "test1@example.com"}},
			Action: func(service *subscription.Service, subscribers []port.User) error {
				return service.Subscribe(&subscribers[0])
			},
			ExpectedError: subscription.ErrAlreadySubscribed,
		},
		{
			Name:        "Get all subscriptions",
			Subscribers: []port.User{},
			Action: func(service *subscription.Service, subscribers []port.User) error {
				_, err := service.Subscriptions()
				return err
			},
			ExpectedResult: []port.User{{Email: "test1@example.com"}},
		},
		{
			Name: "Subscribe multiple emails",
			Subscribers: []port.User{
				{Email: "test2@example.com"},
				{Email: "test3@example.com"},
			},
			Action: func(service *subscription.Service, subscribers []port.User) error {
				for _, subscriber := range subscribers {
					if err := service.Subscribe(&subscriber); err != nil {
						return err
					}
				}
				return nil
			},
			ExpectedResult: []port.User{
				{Email: "test1@example.com"},
				{Email: "test2@example.com"},
				{Email: "test3@example.com"},
			},
		},
		{
			Name: "Subscribe new and already subscribed emails",
			Subscribers: []port.User{
				{Email: "test4@example.com"},
				{Email: "test1@example.com"},
			},
			Action: func(service *subscription.Service, subscribers []port.User) error {
				for _, subscriber := range subscribers {
					err := service.Subscribe(&subscriber)
					if err != nil && !errors.Is(err, subscription.ErrAlreadySubscribed) {
						return err
					}
				}
				return nil
			},
			ExpectedResult: []port.User{
				{Email: "test1@example.com"},
				{Email: "test2@example.com"},
				{Email: "test3@example.com"},
				{Email: "test4@example.com"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			runTest(t, tt, service)
		})
	}
}

func runTest(t *testing.T, test SubscriptionTest, service *subscription.Service) {
	err := test.Action(service, test.Subscribers)
	checkError(t, err, test.ExpectedError)
	checkExpectedResult(t, service, test.ExpectedResult)
}

func checkError(t *testing.T, err error, expectedError error) {
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected error %v, got: %v", expectedError, err)
	}
}

func checkExpectedResult(
	t *testing.T,
	service *subscription.Service,
	expectedResult []port.User,
) {
	if expectedResult == nil {
		return
	}

	subscriptions, err := service.Subscriptions()
	if err != nil {
		t.Fatalf("Failed to get all subscriptions: %v", err)
	}

	if !cmp.Equal(subscriptions, expectedResult) {
		t.Errorf("Unexpected subscriptions. Got: %v, Expected: %v", subscriptions, expectedResult)
	}
}
