package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type CurrencyInterface interface {
	ConvertToThb(ctx context.Context, amount float64) Decimal
}

type CurrencyService struct {
	endpoint string
	ttl      time.Duration
	client   *http.Client

	mu          sync.Mutex
	rate        float64
	fetched     time.Time
	nextAttempt time.Time
}

const (
	fxRetryBackoff   = time.Minute
	fxRequestTimeout = 10 * time.Second
)

func NewCurrencyService(endpoint string, ttl time.Duration) *CurrencyService {
	return &CurrencyService{
		endpoint: endpoint,
		ttl:      ttl,
		client:   &http.Client{Timeout: fxRequestTimeout},
		rate:     DefaultUSDToTHBRate,
	}
}

func (s *CurrencyService) ConvertToThb(ctx context.Context, amount float64) Decimal {
	return ConvertToThb(amount, s.getRate(ctx))
}

func (s *CurrencyService) getRate(ctx context.Context) float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	if time.Now().Before(s.nextAttempt) {
		return s.rate
	}

	if s.rate > 0 && time.Since(s.fetched) < s.ttl {
		return s.rate
	}

	newRate, err := s.fetchRate(ctx)
	if err != nil {
		slog.WarnContext(ctx, "CurrencyService.fetchRate failed, using fallback rate", "error", err, "rate", s.rate)
		s.nextAttempt = time.Now().Add(fxRetryBackoff)
		return s.rate
	}

	s.rate = newRate
	s.fetched = time.Now()
	return s.rate
}

func (s *CurrencyService) fetchRate(ctx context.Context) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.endpoint, nil)
	if err != nil {
		return 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var payload struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, err
	}

	rate, ok := payload.Rates["THB"]
	if !ok || rate <= 0 {
		return 0, fmt.Errorf("THB rate not found in response")
	}
	return rate, nil
}
