package storages

import (
	"fmt"

	"github.com/gogotchuri/gocialite"
)

//In memory storage for Gocialite
var _ gocialite.GocialStorage = &MemoryStorage{}

type MemoryStorage struct {
	storage map[string]*gocialite.Gocial
}

//NewMemoryStorage returns a new MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		storage: make(map[string]*gocialite.Gocial),
	}
}

//Get a Gocialite struct from memory
func (s *MemoryStorage) Get(key string) (*gocialite.Gocial, error) {
	if val, ok := s.storage[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("Key %s not found", key)
}

//Set a Gocialite struct to memory
func (s *MemoryStorage) Set(key string, value *gocialite.Gocial) error {
	s.storage[key] = value
	return nil
}

//Set a Gocialite struct to memory
func (s *MemoryStorage) Delete(key string) error {
	delete(s.storage, key)
	return nil
}
