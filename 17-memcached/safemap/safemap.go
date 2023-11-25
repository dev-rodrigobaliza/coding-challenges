package safemap

import (
	"sync"
)

type SafeMap struct {
	sync.RWMutex
	data map[string][]byte
}

func New() *SafeMap {
	return &SafeMap{
		data: make(map[string][]byte),
	}
}

func (s *SafeMap) Clear() {
	s.Lock()
	s.data = make(map[string][]byte)
	s.Unlock()
}

func (s *SafeMap) Delete(key string) {
	s.Lock()
	delete(s.data, key)
	s.Unlock()
}

func (s *SafeMap) Get(key string) []byte {
	s.RLock()
	defer s.RUnlock()
	
	value, ok := s.data[key]
	if !ok {
		return nil
	}

	return value
}

func (s *SafeMap) Set(key string, value []byte) {
	s.Lock()
	s.data[key] = value
	s.Unlock()
}
