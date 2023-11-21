package safemap

import (
	"errors"
	"sync"
)

var (
	errNotFoundKey = errors.New("key not found")
)

type SafeMap struct {
	sync.RWMutex
	data map[string]string
}

func New() *SafeMap {
	return &SafeMap{
		data: make(map[string]string),
	}
}

func (s *SafeMap) Clear() {
	s.Lock()
	s.data = make(map[string]string)
	s.Unlock()
}

func (s *SafeMap) Delete(key string) {
	s.Lock()
	delete(s.data, key)
	s.Unlock()
}

func (s *SafeMap) Get(key string) string {
	s.RLock()
	value, ok := s.data[key]
	if !ok {
		return ""
	}
	s.RUnlock()

	return value
}

func (s *SafeMap) Set(key string, value string) {
	s.Lock()
	s.data[key] = value
	s.Unlock()
}