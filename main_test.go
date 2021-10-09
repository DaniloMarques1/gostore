package main

import (
	"testing"
)

func TestParseMessage(t *testing.T) {
	msg := "op=store;key=name;value=Danilo"
	message, err := parseMessage(msg)
	assertNil(t, err)
	assertEqual(t, message.key, "name")
	assertEqual(t, message.value, "Danilo")
	assertEqual(t, message.op, "store")
}

func TestParseMessageError1(t *testing.T) {
	msg := "op=store,key=value"
	_, err := parseMessage(msg)
	assertNotNil(t, err)
	assertEqual(t, InvalidSyntax, err.Error())
}

func TestParseMessageError2(t *testing.T) {
	msg := "daniliisnice"
	_, err := parseMessage(msg)
	assertNotNil(t, err)
	assertEqual(t, InvalidSyntax, err.Error())
}

func TestParseMessageError3(t *testing.T) {
	msg := "da;ni;li=isnice"
	_, err := parseMessage(msg)
	assertNotNil(t, err)
	assertEqual(t, InvalidSyntax, err.Error())
}

func TestParseMessageError4(t *testing.T) {
	msg := "op=delete;ni;li=isnice"
	_, err := parseMessage(msg)
	assertNotNil(t, err)
	assertEqual(t, InvalidSyntax, err.Error())
}

func TestParseMessageError5(t *testing.T) {
	msg := "op=store;key=name;value"
	_, err := parseMessage(msg)
	assertNotNil(t, err)
	assertEqual(t, InvalidSyntax, err.Error())
}

func TestParseMessageError6(t *testing.T) {
	msg := "opa=store;opa=name;opa=12"
	_, err := parseMessage(msg)
	assertNotNil(t, err)
	assertEqual(t, InvalidMessageKey, err.Error())
}

func TestGetOpFromMessage(t *testing.T) {
	msg := Message{
		op:    OP_STORE,
		key:   "name",
		value: "Danilo",
	}

	op, err := getOperationFromMessage(&msg)
	assertNil(t, err)
	opType, ok := op.(StoreOperation) // get the concrete type
	assertTrue(t, ok)
	assertEqual(t, opType.value, "Danilo")
}

func TestGetOpFromMessageError1(t *testing.T) {
	msg := Message{
		op:    "sum", // not supported op
		key:   "name",
		value: "Danilo",
	}

	_, err := getOperationFromMessage(&msg)
	assertEqual(t, OperationNotSupported, err.Error())
}

func TestGetOpFromMessageError2(t *testing.T) {
	msg := "op=sum;key=name;value=Danilo"
	message, err := parseMessage(msg)
	assertNil(t, err)

	_, err = getOperationFromMessage(message)
	assertNotNil(t, err)
	assertEqual(t, OperationNotSupported, err.Error())
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Fatalf("\nAssert equals error\n\tExpected %v\n\tActual %v\n", expected, actual)
	}
}

func assertTrue(t *testing.T, value bool) {
	if !value {
		t.Fatalf("\nAssert true error\n\tgot %v\n", value)
	}
}

// only work for err type so far
func assertNil(t *testing.T, err interface{}) {
	if err != nil {
		t.Fatalf("\nAssert nil error\n\tgot %v\n", err)
	}
}

func assertNotNil(t *testing.T, err interface{}) {
	if err == nil {
		t.Fatalf("\nAssert not nil error\n\tgot %v\n", err)
	}
}
