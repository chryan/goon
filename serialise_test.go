package goon

import (
	"testing"
)

var unit1 = &testUnit{
	Name: "testUnit 1",
	Type: new(int64),
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

func TestSerialise(t *testing.T) {
	Marshal(unit1)
}
