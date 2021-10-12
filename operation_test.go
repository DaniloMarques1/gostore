package main

import (
	"testing"
)

func TestOperationStore(t *testing.T) {
	storage := NewStorage()
	storeOp := StoreOperation{key: "name", value: "Danilo"}
	resp, err := storeOp.executeOperation(storage)

	assertNil(t, err)
	assertEqual(t, "Danilo", storage.Read("name"))
	assertEqual(t, StoredSuccessFully, resp)
}

func TestOperationRead1(t *testing.T) {
	storage := NewStorage()
	storage.Store("name", "Danilo")
	readOp := ReadOperation{key: "name"}
	resp, err := readOp.executeOperation(storage)

	assertNil(t, err)
	assertEqual(t, "Danilo", resp)

}

func TestOperationRead2(t *testing.T) {
	storage := NewStorage()
	storage.Store("values", []int{1, 2, 3})
	readOp := ReadOperation{key: "values"}
	resp, err := readOp.executeOperation(storage)

	list := resp.([]int)

	assertNil(t, err)
	assertEqual(t, 3, len(list))
}

func TestOperationDelete(t *testing.T) {
	storage := NewStorage()
	storage.Store("name", "Danilo")
	deleteOp := DeleteOperation{key: "name"}
	resp, err := deleteOp.executeOperation(storage)

	assertNil(t, err)
	assertEqual(t, DeletedSuccessFully, resp)
}

func TestOperationList(t *testing.T) {
	storage := NewStorage()
	storage.Store("name", "Danilo")
	storage.Store("values", []int{1, 2, 3})
	storage.Store("age", 22)

	listOp := ListOperation{}
	resp, err := listOp.executeOperation(storage)

	list := resp.([]interface{})

	assertNil(t, err)
	assertEqual(t, 3, len(list))
	assertEqual(t, "Danilo", list[0])
}
