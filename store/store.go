package store

type Store struct {
	KeyValueStore *KeyValueStore
}

func NewStore() *Store {
	return &Store{
		KeyValueStore: NewKeyValueStore(),
	}
}
