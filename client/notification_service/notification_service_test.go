package notification_service_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	notificationService "bitbucket.org/moladinTech/go-lib-common/client/notification_service"
	sentryMock "bitbucket.org/moladinTech/go-lib-common/sentry/mocks"
	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockSentry struct{}

func (m *MockSentry) StartSpan(context.Context, string) interface{} {
	return nil
}

func (m *MockSentry) Finish(interface{}) {}

func TestNotification_SendSlack(t *testing.T) {
	span := sentry.Span{}
	mockSentry := sentryMock.NewISentry(t)
	mockSentry.On("StartSpan", mock.Anything, mock.Anything).
		Return(&span).
		Once()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := notificationService.Response{
			Status: "success",
			Data:   nil,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonResponse))
	}))
	defer server.Close()

	ns := notificationService.NewNotificationService(
		validator.New(),
		notificationService.WithSentry(mockSentry),
		notificationService.WithServiceEnv("development"),
		notificationService.WithServiceName("moladin-go-order-service"),
		notificationService.WithTimeoutInSeconds(10),
		notificationService.WithURL(server.URL),
	)

	err := ns.SendSlack(context.Background(), "test-notification-channel", "Test Message")
	require.NotNil(t, ns)
	require.NoError(t, err)
}

func TestNotification_SendSlackError(t *testing.T) {
	span := sentry.Span{}
	mockSentry := sentryMock.NewISentry(t)
	mockSentry.On("StartSpan", mock.Anything, mock.Anything).
		Return(&span).
		Once()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := notificationService.Response{
			Status: "fail",
			Data:   nil,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string(jsonResponse))
	}))
	defer server.Close()

	ns := notificationService.NewNotificationService(
		validator.New(),
		notificationService.WithSentry(mockSentry),
		notificationService.WithServiceEnv("development"),
		notificationService.WithServiceName("moladin-go-order-service"),
		notificationService.WithTimeoutInSeconds(10),
		notificationService.WithURL(server.URL),
	)

	err := ns.SendSlack(context.Background(), "test-notification-channel", "Test Message")
	require.NotNil(t, ns)
	require.Error(t, err)
}

//

func TestNotification_SendEmail(t *testing.T) {
	span := sentry.Span{}
	mockSentry := sentryMock.NewISentry(t)
	mockSentry.On("StartSpan", mock.Anything, mock.Anything).
		Return(&span).
		Once()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := notificationService.Response{
			Status: "success",
			Data:   nil,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonResponse))
	}))
	defer server.Close()

	ns := notificationService.NewNotificationService(
		validator.New(),
		notificationService.WithSentry(mockSentry),
		notificationService.WithServiceEnv("development"),
		notificationService.WithServiceName("moladin-go-order-service"),
		notificationService.WithTimeoutInSeconds(10),
		notificationService.WithURL(server.URL),
	)

	err := ns.SendEmail(context.Background(), []string{"employee@moladin.com"}, "Test Subject", "Test Message")
	require.NotNil(t, ns)
	require.NoError(t, err)
}

func TestNotification_SendEmailError(t *testing.T) {
	span := sentry.Span{}
	mockSentry := sentryMock.NewISentry(t)
	mockSentry.On("StartSpan", mock.Anything, mock.Anything).
		Return(&span).
		Once()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := notificationService.Response{
			Status: "fail",
			Data:   nil,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string(jsonResponse))
	}))
	defer server.Close()

	ns := notificationService.NewNotificationService(
		validator.New(),
		notificationService.WithSentry(mockSentry),
		notificationService.WithServiceEnv("development"),
		notificationService.WithServiceName("moladin-go-order-service"),
		notificationService.WithTimeoutInSeconds(10),
		notificationService.WithURL(server.URL),
	)

	err := ns.SendEmail(context.Background(), []string{"employee@moladin.com"}, "Test Subject", "Test Message")
	require.NotNil(t, ns)
	require.Error(t, err)
}
