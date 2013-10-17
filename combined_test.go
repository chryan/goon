package goon

import (
	"reflect"
	"testing"
)

type TestUint uint16

var testMapPtr = map[string]interface{}{
	"Key1": 1,
	"Key2": 2,
	"Key3": 3,
}

func TestGoon(t *testing.T) {
	compareUnit.Ignored = "IGNORE_ME!"
	serialisemap := map[string]interface{}{
		"unit1": compareUnit,
		//"testptr": testUnitMapPtr,
	}

	bytes, err := Marshal(serialisemap, "goon")
	if err != nil {
		t.Fatalf("Serialisation errors: %v", err)
	}
	compareUnit.Ignored = ""

	deserialisemap, errs := UnmarshalTyped("data.goon", bytes, new(TestTypeFactory))
	if errs != nil {
		t.Fatalf("Deserialisation errors: %v", *errs)
	}

	if !reflect.DeepEqual(serialisemap, deserialisemap) {
		t.Fatalf("Serialise and deserialise data do not match:")
	}
}
