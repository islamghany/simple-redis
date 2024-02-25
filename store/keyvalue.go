package store

import "time"

type KeyValueEntry struct {
	Data string
	TTL  time.Time
}

type KeyValueStore struct {
	Store map[string]KeyValueEntry
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		Store: make(map[string]KeyValueEntry),
	}
}

func (kvs *KeyValueStore) Set(key string, value string) {
	kvs.Store[key] = KeyValueEntry{
		Data: value,
	}
}

func (kvs *KeyValueStore) SetWithTTL(key string, value string, ttl time.Duration) {
	kvs.Store[key] = KeyValueEntry{
		Data: value,
		TTL:  time.Now().Add(ttl),
	}
}

func (kvs *KeyValueStore) Get(key string) (string, bool) {
	entry, ok := kvs.Store[key]
	if !ok {
		return "", false
	}
	// Check if the entry has expired and ttl in not zero value
	if entry.TTL != (time.Time{}) && entry.TTL.Before(time.Now()) {
		delete(kvs.Store, key)
		return "", false
	}
	return entry.Data, true
}

func (kvs *KeyValueStore) Delete(key string) {
	delete(kvs.Store, key)
}
