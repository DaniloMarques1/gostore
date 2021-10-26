package main

import (
	"testing"
)

func TestStorageStore(t *testing.T) {
	storage := NewStorage()
	key := "name"
	value := "Danilo"
	storage.Store(key, value)

	assertEqual(t, storage.db[key], value)
}

func TestStorageRead(t *testing.T) {
	storage := NewStorage()
	key := "name"
	value := "Danilo"
	storage.Store(key, value)

	assertEqual(t, storage.Read(key), value)
}

func TestStorageDelete(t *testing.T) {
	storage := NewStorage()
	key := "name"
	value := "Danilo"
	storage.Store(key, value)
	storage.Delete(key)

	assertEqual(t, nil, storage.Read(key))
}

func TestStorageList(t *testing.T) {
	storage := NewStorage()
	storage.Store("name", "Danilo")
	storage.Store("id", 1)
	storage.Store("age", 22)
	list := storage.List()

	assertEqual(t, 3, len(list))
}

func TestStorageKeys(t *testing.T) {
	storage := NewStorage()
	storage.Store("name", "Danilo")
	storage.Store("id", 1)
	storage.Store("age", 22)
	list := storage.List()

	assertEqual(t, 3, len(list))
}
