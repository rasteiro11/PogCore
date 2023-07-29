package kafka

import (
	"sync"
)

type attributes struct {
	sync.RWMutex
	attributes map[string][]string
}

func newAttributes() *attributes {
	return &attributes{
		attributes: make(map[string][]string),
	}
}

func (a *attributes) Add(key, value string) {
	a.Lock()
	defer a.Unlock()
	a.attributes[key] = append(a.attributes[key], value)
}

func (a *attributes) Get(key string) string {
	a.RLock()
	defer a.RUnlock()
	values, ok := a.attributes[key]
	if !ok || len(values) == 0 {
		return ""
	}
	return values[0]
}

func (a *attributes) Lookup(key string) (string, bool) {
	a.RLock()
	defer a.RUnlock()
	values, ok := a.attributes[key]
	if !ok || len(values) == 0 {
		return "", false
	}
	return values[0], true
}

func (a *attributes) Delete(key string) {
	a.Lock()
	defer a.Unlock()
	delete(a.attributes, key)
}

func (a *attributes) Values() map[string][]string {
	a.RLock()
	defer a.RUnlock()
	result := make(map[string][]string, len(a.attributes))
	for k, v := range a.attributes {
		values := make([]string, len(v))
		copy(values, v)
		result[k] = values
	}
	return result
}
