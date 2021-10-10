package main

type Storage struct {
	db map[string]string
}

type StorageInterface interface {
	Store(key string, value string)
	Read(key string) string
	Delete(key string)
}


// returns a pointer to a new storage
func NewStorage() *Storage {
	db := make(map[string]string)
	return &Storage{
		db: db,
	}
}

// stores a value with the specified key inside the in memory db
func (s *Storage) Store(key string, value string) {
	s.db[key] = value
}

// retrives a stored value
func (s *Storage) Read(key string) string {
	return s.db[key]
}

// delete a stored value
func (s *Storage) Delete(key string) {
	delete(s.db, key)
}

