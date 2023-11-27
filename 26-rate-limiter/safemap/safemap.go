package safemap

import (
	"sync"
)

type SafeMap struct {
	sync.RWMutex
	data map[string]int
}

func New() *SafeMap {
	return &SafeMap{
		data: make(map[string]int),
	}
}

func (s *SafeMap) Clear() {
	s.Lock()
	s.data = make(map[string]int)
	s.Unlock()
}

func (s *SafeMap) Delete(key string) {
	s.Lock()
	delete(s.data, key)
	s.Unlock()
}

func (s *SafeMap) Has(key string) bool {
	s.RLock()
	_, ok := s.data[key]
	s.RUnlock()

	return ok
}

func (s *SafeMap) Get(key string) int {
	s.RLock()
	defer s.RUnlock()

	value, ok := s.data[key]
	if !ok {
		var int int
		return int
	}

	return value
}

func (s *SafeMap) Inc(key string, delta int) bool {
	s.Lock()
	defer s.Unlock()

	value, ok := s.data[key]
	if !ok {
		return false
	}

	value += delta
	s.data[key] = value

	return true
}

func (s *SafeMap) IncAll(delta int) {
	s.Lock()
	defer s.Unlock()

	for k, v := range s.data {
		v += delta
		s.data[k] = v
	}
}

func (s *SafeMap) Restore(delta int) {
	s.Lock()
	defer s.Unlock()

	for k := range s.data {
		s.data[k] = delta
	}
}

func (s *SafeMap) Set(key string, value int) {
	s.Lock()
	s.data[key] = value
	s.Unlock()
}

func (s *SafeMap) Len() int {
	s.RLock()
	defer s.RUnlock()

	return len(s.data)
}

func (s *SafeMap) GetAllKeys() []string {
	s.RLock()
	keys := make([]string, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}
	s.RUnlock()

	return keys
}