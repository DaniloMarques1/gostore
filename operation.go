package main

import (
	"errors"
)

// operations supported
const (
	OP_STORE  = "store"
	OP_DELETE = "delete"
	OP_READ   = "read"
	OP_LIST   = "list"
	OP_KEYS   = "keys"
)

type Operation interface {
	executeOperation(storage StorageInterface) (interface{}, error)
}

type StoreOperation struct {
	key   string
	value interface{}
}

type ReadOperation struct {
	key string
}

type DeleteOperation struct {
	key string
}

type ListOperation struct{}

type KeysOperation struct{}

func (sop StoreOperation) executeOperation(storage StorageInterface) (interface{}, error) {
	value := storage.Read(sop.key)
	if value != nil {
		return "", errors.New(DuplicationOfKey)
	}

	storage.Store(sop.key, sop.value)
	return StoredSuccessFully, nil
}

// we do use conn because we do not need send message to client
func (dop DeleteOperation) executeOperation(storage StorageInterface) (interface{}, error) {
	value := storage.Read(dop.key)
	if value == nil {
		return "", errors.New(KeyNotFound)
	}

	storage.Delete(dop.key)
	return DeletedSuccessFully, nil
}

// we write to conn the read result
func (rop ReadOperation) executeOperation(storage StorageInterface) (interface{}, error) {
	value := storage.Read(rop.key)
	if value == nil {
		return "", errors.New(KeyNotFound)
	}

	return value, nil
}

func (lop ListOperation) executeOperation(storage StorageInterface) (interface{}, error) {
	resp := storage.List()
	return resp, nil
}

func (kop KeysOperation) executeOperation(storage StorageInterface) (interface{}, error) {
	resp := storage.Keys()
	return resp, nil
}
