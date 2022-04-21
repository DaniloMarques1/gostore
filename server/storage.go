package server

import (
	c "github.com/ostafen/clover"
	"log"
)

type Storage struct {
	db            *c.DB
	collectionName string
}

// TODO have to handle errors
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
func NewStorage() (*Storage, error) {
	db, err := c.Open("db")
	if err != nil {
		return nil, err
	}
	collectionName := "gostore"
	db.CreateCollection(collectionName)
	return &Storage{
		db: db,
		collectionName: collectionName,
	}, nil
}

// stores a value with the specified key inside the in memory db
func (s *Storage) Store(key string, value interface{}) {
	doc := c.NewDocument()
	doc.Set("key", key)
	doc.Set("value", value)
	_, err := s.db.InsertOne(s.collectionName, doc)
	log.Printf("%v\n", err)
}

// retrives a stored value
func (s *Storage) Read(key string) interface{} {
	docs, _ := s.db.Query(s.collectionName).Where(c.Field("key").Eq(key)).FindAll() // TODO
	if len(docs) == 0 {
		return nil
	}
	doc := docs[0]
	return doc.Get("value")
}

// delete a stored value
func (s *Storage) Delete(key string) {
}

func (s *Storage) List() []interface{} {
	resp := make([]interface{}, 0)
	docs, _ := s.db.Query(s.collectionName).FindAll()
	for _, doc := range docs {
		v := doc.Get("value")
		resp = append(resp, v)
	}
	return resp
}

func (s *Storage) Keys() []interface{} {
	resp := make([]interface{}, 0)
	docs, _ := s.db.Query(s.collectionName).FindAll()
	for _, doc := range docs {
		v := doc.Get("key")
		resp = append(resp, v)
	}
	return resp
}

func (s *Storage) Replace(key string, value interface{}) {
	//s.db[key] = value
}
