package server

import (
	"testing"
)

func TestOperationStore1(t *testing.T) {
	storage := NewStorage()
	storeOp := StoreOperation{key: "name", value: "Danilo"}
	resp, err := storeOp.ExecuteOperation(storage)

	assertNil(t, err)
	assertEqual(t, "Danilo", storage.Read("name"))
	assertEqual(t, StoredSuccessFully, resp)
}

func TestOperationStore2(t *testing.T) {
	storage := NewStorage()
	storeOp := StoreOperation{key: "name", value: "Danilo"}
	resp, err := storeOp.ExecuteOperation(storage)

	assertNil(t, err)
	assertEqual(t, "Danilo", storage.Read("name"))
	assertEqual(t, StoredSuccessFully, resp)

	storeOp = StoreOperation{key: "name", value: "Messi"}
	_, err = storeOp.ExecuteOperation(storage)
	assertNotNil(t, err)
}

func TestOperationRead1(t *testing.T) {
	storage := NewStorage()
	storage.Store("name", "Danilo")
	readOp := ReadOperation{key: "name"}
	resp, err := readOp.ExecuteOperation(storage)

	assertNil(t, err)
	assertEqual(t, "Danilo", resp)

}

func TestOperationRead2(t *testing.T) {
	storage := NewStorage()
	storage.Store("values", []int{1, 2, 3})
	readOp := ReadOperation{key: "values"}
	resp, err := readOp.ExecuteOperation(storage)

	list := resp.([]int)

	assertNil(t, err)
	assertEqual(t, 3, len(list))
}

func TestOperationRead3(t *testing.T) {
	storage := NewStorage()
	readOp := ReadOperation{key: "values"}
	_, err := readOp.ExecuteOperation(storage)
	assertNotNil(t, err)
}

func TestOperationDelete1(t *testing.T) {
	storage := NewStorage()
	storage.Store("name", "Danilo")
	deleteOp := DeleteOperation{key: "name"}
	resp, err := deleteOp.ExecuteOperation(storage)

	assertNil(t, err)
	assertEqual(t, DeletedSuccessFully, resp)
}

func TestOperationDelete2(t *testing.T) {
	storage := NewStorage()
	deleteOp := DeleteOperation{key: "name"}
	_, err := deleteOp.ExecuteOperation(storage)
	assertNotNil(t, err)

}

func TestOperationList(t *testing.T) {
	storage := NewStorage()
	storage.Store("name", "Danilo")
	storage.Store("values", []int{1, 2, 3})
	storage.Store("age", 22)

	listOp := ListOperation{}
	resp, err := listOp.ExecuteOperation(storage)

	list := resp.([]interface{})

	assertNil(t, err)
	assertEqual(t, 3, len(list))
}

func TestOperationKeys(t *testing.T) {
	storage := NewStorage()
	storage.Store("name", "Danilo")
	storage.Store("values", []int{1, 2, 3})
	storage.Store("age", 22)

	keysOp := KeysOperation{}
	resp, err := keysOp.ExecuteOperation(storage)
	keys := resp.([]interface{})

	assertNil(t, err)
	assertEqual(t, 3, len(keys))
}
