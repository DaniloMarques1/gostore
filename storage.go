package main

type Storage struct {
	db map[string]interface{}
}

type StorageInterface interface {
	Store(key string, value interface{})
	Read(key string) interface{}
	Delete(key string)
	List() []interface{}
	Keys() []interface{}
	Replace(key string, value interface{})
	//Push(key string, value interface{}) error TODO maybe later
}

// returns a pointer to a new storage
func NewStorage() *Storage {
	db := make(map[string]interface{})
	return &Storage{
		db: db,
	}
}

// stores a value with the specified key inside the in memory db
func (s *Storage) Store(key string, value interface{}) {
	s.db[key] = value
}

// retrives a stored value
func (s *Storage) Read(key string) interface{} {
	return s.db[key]
}

// delete a stored value
func (s *Storage) Delete(key string) {
	delete(s.db, key)
}

func (s *Storage) List() []interface{} {
	resp := make([]interface{}, 0, len(s.db))
	for _, v := range s.db {
		resp = append(resp, v)
	}

	return resp
}

func (s *Storage) Keys() []interface{} {
	resp := make([]interface{}, 0, len(s.db))
	for key := range s.db {
		resp = append(resp, key)
	}

	return resp
}

func (s *Storage) Replace(key string, value interface{}) {
	s.db[key] = value
}
