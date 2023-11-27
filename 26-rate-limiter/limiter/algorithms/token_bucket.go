package algorithms

import (
	"rl/safemap"
	"rl/store"
	"time"
)

type TokenBucket struct {
	limit    int
	buckets  store.Store
	ticker   *time.Ticker
	doneChan chan bool
}

func NewTokenBucket(dsn string, limit int) *TokenBucket {
	buckets := getStore(dsn)
	ticker := time.NewTicker(time.Second)
	done := make(chan bool)

	tk := TokenBucket{
		limit:    limit,
		buckets:  buckets,
		ticker:   ticker,
		doneChan: done,
	}

	go func() {
		for {
			select {
			case <-done:
				return

			case <-ticker.C:
				tk.Inc()
			}
		}
	}()

	return &tk
}

func (tk *TokenBucket) Stop() {
	tk.ticker.Stop()
	tk.doneChan <- true
}

func (tk *TokenBucket) Inc() {
	tk.buckets.IncAll(1)
}

func (tk *TokenBucket) IsAllowed(key string) bool {
	if !tk.buckets.Has(key) {
		tk.buckets.Set(key, tk.limit)
	}

	v := tk.buckets.Get(key)
	if v < 1 {
		return false
	}

	return tk.buckets.Inc(key, -1)
}

func getStore(dsn string) store.Store {
	switch dsn {
	default:
		return safemap.New()
	}
}
