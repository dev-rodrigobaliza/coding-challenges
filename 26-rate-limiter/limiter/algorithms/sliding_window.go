package algorithms

import (
	"container/list"
	"rl/store"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type entry struct {
	id        string
	timestamp int64
}

type SlidingWindow struct {
	sync.Mutex
	limit  int
	dur    time.Duration
	next   uint64
	store  store.Store
	queues map[string]*list.List
}

func NewSlidingWindow(store store.Store, limit int, dur time.Duration) *SlidingWindow {
	tk := SlidingWindow{
		limit:  limit,
		dur:    dur,
		store:  store,
		queues: make(map[string]*list.List),
	}

	return &tk
}

func (sw *SlidingWindow) Stop() {}

func (sw *SlidingWindow) Restore() {}

func (sw *SlidingWindow) IsAllowed(id, key string) bool {
	sw.Lock()
	defer sw.Unlock()

	idx := sw.getIndex(key)
	queue := sw.getQueue(idx)
	if queue == nil {
		// I really hope we never enter here, for gosh sakes...
		return false
	}

	t := time.Now().UnixMicro()
	sw.insertQueue(queue, id, t)
	if queue.Len() == 1 {
		return true
	}

	return sw.checkQueue(id, queue)
}

func (sw *SlidingWindow) getIndex(key string) string {
	if sw.store.Has(key) {
		idx := sw.store.Get(key)
		return strconv.Itoa(idx)
	}

	next := int(atomic.AddUint64(&sw.next, 1))
	sw.store.Set(key, next)
	idx := strconv.Itoa(next)

	sw.queues[idx] = list.New()

	return idx
}

func (sw *SlidingWindow) getQueue(index string) *list.List {
	queue, ok := sw.queues[index]
	if !ok {
		return nil
	}

	return queue
}

func (sw *SlidingWindow) insertQueue(queue *list.List, id string, timestamp int64) {
	e := entry{
		id:        id,
		timestamp: timestamp,
	}
	queue.PushBack(e)
}

func (sw *SlidingWindow) checkQueue(id string, queue *list.List) bool {
	for {
		first := queue.Front()
		if first == nil {
			return true
		}

		e, ok := first.Value.(entry)
		if !ok {
			// I really hope we never enter here, for gosh sakes...
			return false
		}

		w := time.UnixMicro(e.timestamp)
		d := time.Since(w)
		if d <= sw.dur {
			break
		}

		queue.Remove(first)
	}

	return queue.Len() <= sw.limit
}
