package goon

import (
	"reflect"
	"strings"
	"testing"
)

func TestDeserialise(t *testing.T) {
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
		t.Fatalf("Failed to deserialise. Values do not match:\n%+v\n%+v", compareUnit, unit)
	}
}
