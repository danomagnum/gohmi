package hmitags

import (
	"sync"
)

type TagStore struct {
	DriverName string
	m          sync.RWMutex
	data       map[string]any
}

func NewTagStore(name string) *TagStore {

	return &TagStore{
		DriverName: name,
		data:       make(map[string]any),
	}
}

func (ts *TagStore) Read(key string) (any, error) {
	ts.m.RLock()
	defer ts.m.RUnlock()
	val, ok := ts.data[key]
	if !ok {
		return nil, ErrTagNotFound{key: key}
	}
	return val, nil
}

func (ts *TagStore) Write(key string, value any) error {
	ts.m.Lock()
	defer ts.m.Unlock()
	ts.data[key] = value
	return nil
}

func (ts *TagStore) Start() error {
	return nil
}
func (ts *TagStore) Stop() error {
	return nil
}
func (ts *TagStore) Status() string {
	return "ok"
}
func (ts *TagStore) Name() string {
	return ts.DriverName
}
