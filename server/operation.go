package server

import (
	"errors"
)

// operations supported
const (
	OP_STORE   = "store"
	OP_DELETE  = "delete"
	OP_READ    = "read"
	OP_LIST    = "list"
	OP_KEYS    = "keys"
	OP_REPLACE = "replace"
)

// Operation Response Messages
const (
	StoredSuccessFully   = "Value stored successfully"
	DeletedSuccessFully  = "Value removed successfully"
	ReplacedSuccessFully = "Value replaced successfully"
)

type Operation interface {
	ExecuteOperation(storage StorageInterface) (interface{}, error)
	GetOpType() string
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

type ReplaceOperation struct {
	key   string
	value interface{}
}

func (sop StoreOperation) ExecuteOperation(storage StorageInterface) (interface{}, error) {
	value, _ := storage.Read(sop.key) // ignoring error becase we care only if key the exists
	if value != nil {
		return "", errors.New(DuplicationOfKey)
	}

	storage.Store(sop.key, sop.value)
	return StoredSuccessFully, nil
}

func (sop StoreOperation) GetOpType() string {
	return OP_STORE
}

func (dop DeleteOperation) ExecuteOperation(storage StorageInterface) (interface{}, error) {
	_, err := storage.Read(dop.key)
	if err != nil {
		return nil, err
	}

	err = storage.Delete(dop.key)
	if err != nil {
		return nil, err
	}
	return DeletedSuccessFully, nil
}

func (sop DeleteOperation) GetOpType() string {
	return OP_DELETE
}

func (rop ReadOperation) ExecuteOperation(storage StorageInterface) (interface{}, error) {
	value, err := storage.Read(rop.key)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (sop ReadOperation) GetOpType() string {
	return OP_READ
}

func (lop ListOperation) ExecuteOperation(storage StorageInterface) (interface{}, error) {
	resp, err := storage.List()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sop ListOperation) GetOpType() string {
	return OP_LIST
}

func (kop KeysOperation) ExecuteOperation(storage StorageInterface) (interface{}, error) {
	resp, err := storage.Keys()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (sop KeysOperation) GetOpType() string {
	return OP_KEYS
}

func (rop ReplaceOperation) ExecuteOperation(storage StorageInterface) (interface{}, error) {
	_, err := storage.Read(rop.key)
	if err != nil {
		return nil, err
	}
	storage.Replace(rop.key, rop.value)
	return ReplacedSuccessFully, nil
}

func (rop ReplaceOperation) GetOpType() string {
	return OP_REPLACE
}
