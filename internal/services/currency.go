package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"flypro/internal/config"

	"github.com/redis/go-redis/v9"
)

type CurrencyService interface {
	ConvertCurrency(ctx context.Context, amount float64, from, to string) (float64, error)
}

type currencyService struct {
	cfg   *config.Config
	cache *redis.Client
	hc    *http.Client
}

func NewCurrencyService(cfg *config.Config, cache *redis.Client) CurrencyService {
	return &currencyService{
		cfg:   cfg,
		cache: cache,
		hc:    &http.Client{Timeout: 5 * time.Second},
	}
}

// convert amount from->to using exchangerate.host latest endpoint (or similar) and cache rate for 6 hours.
func (s *currencyService) ConvertCurrency(ctx context.Context, amount float64, from, to string) (float64, error) {
	if from == to {
		return amount, nil
	}
	key := fmt.Sprintf("fx:%s:%s", from, to)
	if val, err := s.cache.Get(ctx, key).Float64(); err == nil {
		return amount * val, nil
	}

	u, _ := url.Parse(s.cfg.FXAPIURL)
	q := u.Query()
	q.Set("base", from)
	q.Set("symbols", to)
	if s.cfg.FXAPIKey != "" {
		q.Set("api_key", s.cfg.FXAPIKey)
	}
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	resp, err := s.hc.Do(req)
	if err != nil {
		// try stale cache
		if val, err2 := s.cache.Get(ctx, key).Float64(); err2 == nil {
			return amount * val, nil
		}
		return 0, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	var data struct {
		Rates map[string]float64 `json:"rates"`
	}
	_ = json.Unmarshal(b, &data)
	rate := data.Rates[to]
	if rate <= 0 {
		// try reverse rate
		key2 := fmt.Sprintf("fx:%s:%s", to, from)
		if val, err2 := s.cache.Get(ctx, key2).Float64(); err2 == nil && val > 0 {
			if val != 0 {
				rate = 1.0 / val
			}
		}
	}
	if rate <= 0 {
		rate = 1
	} // fallback safe-guard

	_ = s.cache.Set(ctx, key, rate, 6*time.Hour).Err()
	return amount * rate, nil
}
