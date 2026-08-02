package common_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"store-service/internal/common"

	"github.com/stretchr/testify/assert"
)

func Test_CurrencyService_ConvertToThb_Should_Fetch_Rate_From_API(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"amount":1.0,"base":"USD","date":"2026-01-01","rates":{"THB":36.0}}`))
	}))
	defer server.Close()

	service := common.NewCurrencyService(server.URL, time.Hour)
	actual := service.ConvertToThb(context.Background(), 100)

	assert.Equal(t, common.ConvertToThb(100, 36.0), actual)
}

func Test_CurrencyService_ConvertToThb_Should_Cache_Rate_Until_TTL(t *testing.T) {
	hits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"rates":{"THB":36.0}}`))
	}))
	defer server.Close()

	service := common.NewCurrencyService(server.URL, time.Hour)
	service.ConvertToThb(context.Background(), 100)
	service.ConvertToThb(context.Background(), 200)
	service.ConvertToThb(context.Background(), 300)

	assert.Equal(t, 1, hits)
}

func Test_CurrencyService_ConvertToThb_Should_Refetch_After_TTL_Expires(t *testing.T) {
	hits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"rates":{"THB":36.0}}`))
	}))
	defer server.Close()

	service := common.NewCurrencyService(server.URL, 5*time.Millisecond)
	service.ConvertToThb(context.Background(), 100)
	time.Sleep(20 * time.Millisecond)
	service.ConvertToThb(context.Background(), 200)

	assert.Equal(t, 2, hits)
}

func Test_CurrencyService_ConvertToThb_Should_Fall_Back_To_Default_Rate_On_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	service := common.NewCurrencyService(server.URL, time.Hour)
	actual := service.ConvertToThb(context.Background(), 123)

	assert.Equal(t, common.ConvertToThb(123, common.DefaultUSDToTHBRate), actual)
}
