package goon

import (
	"testing"
	"fmt"
)

var unit1 = &testUnit{
	Name: "testUnit 1",
	Type: nil,
	Map: map[string]int{
		"Key1": 10,
		"Key2": 20,
	},
	Seq: []interface{}{
		0,
		&testUnit{
			Name: "Test",
			Type: new(int64),
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

var testMapPtr = &map[string]interface{}{
	"Key1": 1,
	"Key2": 2,
	"Key3": 3,
}

func TestSerialise(t *testing.T) {
	fmt.Println(unit1.Map)
	m := map[string]interface{}{
		"unit1": unit1,
		"intvar": 10,
		"testMap": testMapPtr,
	}
	if bytes, err := Marshal(m, "datagoon"); err == nil {
		fmt.Println(string(bytes))
	}
}
