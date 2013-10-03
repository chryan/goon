package goon

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
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
	TestInnerValue float32
}

type testUnit struct {
	Name            string
	Type            *int64
	Map             map[string]int
	Seq             []interface{}
	Bool            bool
	TestInner       testInner
	Array           []int
	InterfaceStruct interface{}
	InterfaceVal    interface{}
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
	TestInner: &testInner{
		TestInnerValue: 20.0,
	},
	Array: []int{
		1,
		2,
		3,
		4,
		"GASD",
	},
	InterfaceStruct: &testInner{
		TestInnerValue: 30.0,
	},
	InterfaceVal: "String",
}`)

func TestDeserialise(t *testing.T) {
	var unitType int64
	unitType = 10
	compareUnit := &testUnit{
		Name: "testUnit 1",
		Type: &unitType,
		Map: map[string]int{
			"Key1": 10,
		},
		Seq: []interface{}{
			0,
			&testUnit{
				Name: "Test",
				Type: &unitType,
				Bool: true,
			},
			2,
			3,
		},
		TestInner: testInner{
			TestInnerValue: 20.0,
		},
		Array: []int{
			1,
			2,
			3,
			4,
		},
		InterfaceStruct: &testInner{
			TestInnerValue: 30.0,
		},
		InterfaceVal: "String",
	}

	valuemap, errs := UnmarshalTyped("data.goon", complexTypeTest, new(TestTypeFactory))

	if errs != nil {
		t.Logf("\n%v", strings.Join(errs.Msgs, "\n"))
	}

	val, ok := valuemap["unit1"]
	if !ok {
		t.Fatalf("Failed to deserialise unit1.")
	}
	unit, ok := val.(*testUnit)
	if !ok {
		t.Fatalf("Failed to deserialise unit1 with type Unit")
	}

	if !reflect.DeepEqual(compareUnit, unit) {
		cuj, _ := json.Marshal(compareUnit)
		uj, _ := json.Marshal(unit)
		t.Fatalf("Deserialised values do not match:\n%v\n%v", string(cuj), string(uj))
	}
}

func TestDeserialiseNoType(t *testing.T) {
	//valuemap, errs := Unmarshal("data.goon", complexTypeTest)
	//fmt.Println(valuemap, errs)
}
