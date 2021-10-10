package main

import (
	"errors"
	"log"
)

type Operation interface {
	executeOperation(storage StorageInterface) (string, error)
}

type StoreOperation struct {
	key   string
	value string
}

type ReadOperation struct {
	key string
}

type DeleteOperation struct {
	key string
}

func (sop StoreOperation) executeOperation(storage StorageInterface) (string, error) {
	log.Printf("Executing store operation\n")
	log.Printf("KEY = %v - VALUE = %v\n", sop.key, sop.value)
	value := storage.Read(sop.key)
	if len(value) != 0 {
		return "", errors.New(DuplicationOfKey)
	}

	storage.Store(sop.key, sop.value)
	return StoredSuccessFully, nil
}

// we do use conn because we do not need send message to client
func (dop DeleteOperation) executeOperation(storage StorageInterface) (string, error) {
	log.Printf("Executing delete operation\n")
	log.Printf("KEY = %v\n", dop.key)
	value := storage.Read(dop.key)
	if len(value) == 0 {
		return "", errors.New(KeyNotFound)
	}

	storage.Delete(dop.key)
	return DeletedSuccessFully, nil
}

// we write to conn the read result
func (rop ReadOperation) executeOperation(storage StorageInterface) (string, error) {
	log.Printf("Executing read operation\n")
	log.Printf("KEY = %v\n", rop.key)
	value := storage.Read(rop.key)
	if len(value) == 0 {
		return "", errors.New(KeyNotFound)
	}

	return value, nil
}

