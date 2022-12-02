package util

import "sync"

type SafeMap struct {
	l sync.RWMutex
	M map[string]interface{}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		l: sync.RWMutex{},
		M: make(map[string]interface{}),
	}
}

func (s *SafeMap) Get(key string) (interface{}, bool) {
	s.l.RLock()
	defer s.l.RUnlock()
	v, ok := s.M[key]
	return v, ok
}

func (s *SafeMap) Set(key string, value interface{}) {
	s.l.Lock()
	s.M[key] = value
	s.l.Unlock()
}

func (s *SafeMap) Len() int {
	s.l.RLock()
	defer s.l.RUnlock()
	return len(s.M)
}

func (s *SafeMap) Del(key string) {
	s.l.Lock()
	delete(s.M, key)
	s.l.Unlock()
}
