package algorithms

import (
	"rl/store"
	"time"
)

type FixedWindow struct {
	limit    int
	store    store.Store
	ticker   *time.Ticker
	doneChan chan bool
}

func NewFixedWindow(store store.Store, limit int, dur time.Duration) *FixedWindow {
	ticker := time.NewTicker(dur)
	done := make(chan bool)

	tk := FixedWindow{
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
				tk.Restore()
			}
		}
	}()

	return &tk
}

func (fw *FixedWindow) Stop() {
	fw.ticker.Stop()
	fw.doneChan <- true
}

func (fw *FixedWindow) Restore() {
	fw.store.Restore(fw.limit)
}

func (fw *FixedWindow) IsAllowed(key string) bool {
	if !fw.store.Has(key) {
		fw.store.Set(key, fw.limit)
	}

	v := fw.store.Get(key)
	if v < 1 {
		return false
	}

	return fw.store.Inc(key, -1)
}
