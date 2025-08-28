package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"flypro/internal/config"

	"github.com/redis/go-redis/v9"
)

// CurrencyService defines the behavior for currency conversion
type CurrencyService interface {
	ConvertCurrency(ctx context.Context, amount float64, from, to string) (float64, error)
}

type currencyService struct {
	cfg   *config.Config
	cache *redis.Client
	hc    *http.Client
}

// NewCurrencyService creates a new currency service
func NewCurrencyService(cfg *config.Config, cache *redis.Client) CurrencyService {
	return &currencyService{
		cfg:   cfg,
		cache: cache,
		hc:    &http.Client{Timeout: 5 * time.Second},
	}
}

// API response structure from exchangerate-api.com
type exchangeRateResponse struct {
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

// ConvertCurrency converts amount from one currency to another using exchangerate-api.com
func (s *currencyService) ConvertCurrency(ctx context.Context, amount float64, from, to string) (float64, error) {
	if from == to {
		return amount, nil
	}

	cacheKey := fmt.Sprintf("fx:%s:%s", from, to)
	if val, err := s.cache.Get(ctx, cacheKey).Float64(); err == nil {
		return amount * val, nil
	}

	// Build API URL: https://v6.exchangerate-api.com/v6/{API_KEY}/latest/{BASE}
	apiURL := fmt.Sprintf("%s/%s/latest/%s", s.cfg.FXAPIURL, s.cfg.FXAPIKey, from)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	resp, err := s.hc.Do(req)
	if err != nil {
		// fallback to stale cache
		if val, err2 := s.cache.Get(ctx, cacheKey).Float64(); err2 == nil {
			return amount * val, nil
		}
		return 0, fmt.Errorf("fx api request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data exchangeRateResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("failed to parse fx api response: %w", err)
	}

	rate := data.ConversionRates[to]
	if rate <= 0 {
		return 0, fmt.Errorf("invalid rate for %s -> %s", from, to)
	}

	// cache rate for 6 hours
	_ = s.cache.Set(ctx, cacheKey, rate, 6*time.Hour).Err()

	return amount * rate, nil
}
