package safemap

import (
	"sync"
)

type SafeMap[D any] struct {
	sync.RWMutex
	data map[string]D
}

func New[D any]() *SafeMap[D] {
	return &SafeMap[D]{
		data: make(map[string]D),
	}
}

func (s *SafeMap[D]) Clear() {
	s.Lock()
	s.data = make(map[string]D)
	s.Unlock()
}

func (s *SafeMap[D]) Delete(key string) {
	s.Lock()
	delete(s.data, key)
	s.Unlock()
}

func (s *SafeMap[D]) Get(key string) D {
	s.RLock()
	defer s.RUnlock()

	value, ok := s.data[key]
	if !ok {
		var d D
		return d
	}

	return value
}

func (s *SafeMap[D]) Set(key string, value D) {
	s.Lock()
	s.data[key] = value
	s.Unlock()
}

func (s *SafeMap[D]) Len() int {
	s.RLock()
	defer s.RUnlock()

	return len(s.data)
}