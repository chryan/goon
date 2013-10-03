package goon

import (
	"testing"
	"strings"
	"fmt"
	//"reflect"
)

type TestTypeFactory struct {

}

func (t *TestTypeFactory) New(typename, pkgname string) interface{} {
	switch typename {
	case "testInner":
		return new(testInner)
	case "testUnit":
		return new(testUnit)
	}
	return nil
}

type testInner struct {
	testInnerValue float32
}

type testUnit struct {
	Name  string
	Type  *int64
	Map   map[string]int
	Seq	  []interface{}
	Bool  bool
	testInner testInner
	Array []int
	InterfaceStruct interface{}
	InterfaceVal interface{}
}

var complexTypeTest []byte = []byte(`package testdata

var unit1 = &testUnit{
	Name: "testUnit 1",
	Type: 10,
	Map: map[string]int{
		Key1: 10,
		Key2: "Error",
	},
	Seq: []interface{}{
		0,
		&testUnit{
			Name: "Test",
			Type: 10,
			Bool: true,
		},
		2,
		3,
	},
	testInner: &testInner{
		testInnerValue: 20.0,
	},
	Array: []int{
		1,
		2,
		3,
		4,
		"GASD",
	},
	InterfaceStruct: &testInner{
		testInnerValue: 30.0,
	},
	InterfaceVal: "String",
}`)

func TestDeserialise(t *testing.T) {
	valuemap, errs := UnmarshalTyped("data.goon", complexTypeTest, new(TestTypeFactory))

	if errs != nil {
		t.Logf("\n%v", strings.Join(errs.Msgs, "\n"))
	}

	var val interface{}
	var ok bool

	val, ok = valuemap["unit1"]
	if !ok {
		t.Fatalf("Failed to deserialise unit1.")
	}

	unit, ok := val.(*testUnit)
	if !ok {
		t.Fatalf("unit1 not of type testUnit.")
	}

	switch {
	case unit.Name != "testUnit 1":
		t.Fatalf("Failed to deserialise unit1.Name")
	case *unit.Type != 10:
		t.Fatalf("Failed to deserialise unit1.Type")
	}
}

func TestDeserialiseNoType(t *testing.T) {
	valuemap, errs := Unmarshal("data.goon", complexTypeTest)
	fmt.Println(valuemap, errs)
}