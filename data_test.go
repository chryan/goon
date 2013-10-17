package goon

import (
	"go/token"
)

type testInner struct {
	TestInnerValue float32
}

type testUnit struct {
	Name            string
	Type            *int64
	Map             map[string]int
	Uint            uint
	InterfaceMap    map[int]interface{}
	Seq             []interface{}
	Bool            bool
	TestInner       testInner
	Array           []int
	InterfaceStruct interface{}
	InterfaceVal    interface{}
	Position        token.Position
	Ignored         string `goon:"ignore"`
}

type TestTypeFactory struct {
}

func (t *TestTypeFactory) New(typename, pkgname string) interface{} {
	switch typename {
	case "testInner":
		return new(testInner)
	case "testUnit":
		return new(testUnit)
	}

	switch {
	case pkgname == "token" && typename == "Position":
		return new(token.Position)
	}
	return nil
}

var complexTypeTest []byte = []byte(`package goon

var unit1 = &testUnit{
	Name: "testUnit 1",
	Type: 10,
	Map: map[string]int{
		"Key1": 10,
	},
	Uint: 23,
	InterfaceMap: map[int]interface{}{
		1: "Test1",
		2: 2,
		3: "Test3",
		4: &testUnit{
			Name: "Test",
			Type: 10,
			Map: nil,
			Uint: 0,
			InterfaceMap: nil,
			Seq: nil,
			Bool: true,
			TestInner: testInner{
				TestInnerValue: 0,
			},
			Array: nil,
			InterfaceStruct: nil,
			InterfaceVal: nil,
			Position: token.Position{
				Filename: "",
				Offset: 0,
				Line: 0,
				Column: 0,
			},
		},
	},
	Seq: []interface{}{
		0,
		&testUnit{
			Name: "Test",
			Type: 10,
			Map: nil,
			Uint: 0,
			InterfaceMap: nil,
			Seq: nil,
			Bool: true,
			TestInner: testInner{
				TestInnerValue: 0,
			},
			Array: nil,
			InterfaceStruct: nil,
			InterfaceVal: nil,
			Position: token.Position{
				Filename: "",
				Offset: 0,
				Line: 0,
				Column: 0,
			},
		},
		2,
		3,
	},
	Bool: false,
	TestInner: testInner{
		TestInnerValue: 20,
	},
	Array: []interface{}{
		1,
		2,
		3,
		4,
	},
	InterfaceStruct: &testInner{
		TestInnerValue: 30,
	},
	InterfaceVal: "String",
	Position: token.Position{
		Filename: "Word",
		Offset: 0,
		Line: 0,
		Column: 0,
	},
}`)

var unitType int64 = 10

var compareUnit = &testUnit{
	Name: "testUnit 1",
	Type: &unitType,
	Map: map[string]int{
		"Key1": 10,
	},
	Uint: 23,
	InterfaceMap: map[int]interface{}{
		1: "Test1",
		2: int64(2),
		3: "Test3",
		4: &testUnit{
			Name: "Test",
			Type: &unitType,
			Bool: true,
		},
	},
	Seq: []interface{}{
		int64(0),
		&testUnit{
			Name: "Test",
			Type: &unitType,
			Bool: true,
		},
		int64(2),
		int64(3),
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
	Position: token.Position{
		Filename: "Word",
	},
}
