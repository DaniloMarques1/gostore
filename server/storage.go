package server

import (
	c "github.com/ostafen/clover"
	"errors"
)

type Storage struct {
	db             *c.DB
	collectionName string
}

type StorageInterface interface {
	Store(key string, value interface{}) error
	Read(key string) (interface{}, error)
	Delete(key string) error
	List() ([]interface{}, error)
	Keys() ([]interface{}, error)
	Replace(key string, value interface{}) error
	//Push(key string, value interface{}) error TODO maybe later
}

// returns a pointer to a new storage
func NewStorage() (*Storage, error) {
	db, err := c.Open("db")
	if err != nil {
		return nil, err
	}
	collectionName := "gostore"
	db.CreateCollection(collectionName)
	return &Storage{
		db:             db,
		collectionName: collectionName,
	}, nil
}

// stores a value with the specified key inside the in memory db
func (s *Storage) Store(key string, value interface{}) error {
	doc := c.NewDocument()
	doc.Set("key", key)
	doc.Set("value", value)
	_, err := s.db.InsertOne(s.collectionName, doc)
	return err
}

// retrives a stored value
func (s *Storage) Read(key string) (interface{}, error) {
	docs, _ := s.db.Query(s.collectionName).Where(c.Field("key").Eq(key)).FindAll()
	if len(docs) == 0 {
		return nil, errors.New(KeyNotFound)
	}
	doc := docs[0]
	return doc.Get("value"), nil
}

// delete a stored value
func (s *Storage) Delete(key string) error {
	return s.db.Query(s.collectionName).Where(c.Field("key").Eq(key)).Delete()
}

func (s *Storage) List() ([]interface{}, error) {
	resp := make([]interface{}, 0)
	docs, err := s.db.Query(s.collectionName).FindAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docs {
		v := doc.Get("value")
		resp = append(resp, v)
	}
	return resp, nil
}

func (s *Storage) Keys() ([]interface{}, error) {
	resp := make([]interface{}, 0)
	docs, err := s.db.Query(s.collectionName).FindAll()
	if err != nil {
		return nil, err
	}
	for _, doc := range docs {
		v := doc.Get("key")
		resp = append(resp, v)
	}
	return resp, nil
}

func (s *Storage) Replace(key string, value interface{}) error {
	update := map[string]interface{}{"key": key, "value": value}
	return s.db.Query(s.collectionName).Where(c.Field("key").Eq(key)).Update(update) // TODO
}
