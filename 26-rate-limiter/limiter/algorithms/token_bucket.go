package algorithms

import (
	"rl/store"
	"time"
)

type TokenBucket struct {
	limit    int
	store    store.Store
	ticker   *time.Ticker
	doneChan chan bool
}

func NewTokenBucket(store store.Store, limit int, dur time.Duration) *TokenBucket {
	ticker := time.NewTicker(dur)
	done := make(chan bool)

	tk := TokenBucket{
		limit:    limit,
		store:    store,
		ticker:   ticker,
		doneChan: done,
	}

	go func() {
		for {
			select {
			case <-done:
				return

			case <-ticker.C:
				tk.IncAll()
			}
		}
	}()

	return &tk
}

func (tk *TokenBucket) Stop() {
	tk.ticker.Stop()
	tk.doneChan <- true
}

func (tk *TokenBucket) IncAll() {
	tk.store.IncAll(1)
}

func (tk *TokenBucket) IsAllowed(key string) bool {
	if !tk.store.Has(key) {
		tk.store.Set(key, tk.limit)
	}

	v := tk.store.Get(key)
	if v < 1 {
		return false
	}

	return tk.store.Inc(key, -1)
}
