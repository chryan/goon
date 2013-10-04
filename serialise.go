package goon

import (
	"reflect"
	"fmt"
	"bytes"
)

type serialiser struct {
	buffer    *bytes.Buffer
	varbuff   *bytes.Buffer

	indent    int
	indentStr string
}

func (s *serialiser) startVar(name string) {
	s.buffer.WriteString(fmt.Sprintf("var %v ", name))
}

func (s *serialiser) startSerialise(v interface{}) {
	s.serialise(v)

	// Consume variable buffer.
	if s.varbuff.Len() > 0 {
		s.buffer.WriteString("= ")
		s.buffer.ReadFrom(s.varbuff)
		s.varbuff.Reset()
	} else {
		s.buffer.WriteString("interface{}")
	}

	s.buffer.WriteString("\n\n")
}

func (s *serialiser) init(pkgname string) {
	s.buffer.WriteString(fmt.Sprintf("package %v\n\n", pkgname))
}

func (s *serialiser) writeStruct(vval reflect.Value, vtype reflect.Type) {
	s.varbuff.WriteString(vval.Type().Name())
	s.varbuff.WriteString("{")
	if numfield := vtype.NumField(); numfield > 0 {
		s.incIndent()
		for i := 0; i < numfield; i++ {
			ftype := vtype.Field(i)
			fval := vval.Field(i)

			s.varbuff.WriteString("\n")
			s.writeIndents()
			s.varbuff.WriteString(fmt.Sprintf("%v: ", ftype.Name))
			s.serialise(fval.Interface())
			s.varbuff.WriteString(",")
		}
		s.varbuff.WriteString("\n")
		s.decIndent()
	}
	s.writeIndents()
	s.varbuff.WriteString("}")
}

func (s *serialiser) writeMap(vval reflect.Value, vtype reflect.Type) {
	fmt.Println(vtype.Key())
}

func (s *serialiser) writeIndents() {
	for i := 0; i < s.indent; i++ {
		s.varbuff.WriteString(s.indentStr)
	}
}

func (s *serialiser) incIndent() {
	s.indent++
}

func (s *serialiser) decIndent() {
	s.indent--
}

func (s *serialiser) writePtr(ptr bool) {
	if ptr {
		s.varbuff.WriteString("&")
	}
}

func (s *serialiser) writeNil() {
	s.varbuff.WriteString("nil")
}

func (s *serialiser) serialise(v interface{}) {
	vval := reflect.ValueOf(v)
	vkind := vval.Kind()

	//oldval := vval

	var vtype reflect.Type
	ptr := false

	if vkind == reflect.Ptr {
		if vval.IsNil() {
			s.writeNil()
			return
		}

		vval = vval.Elem()
		vtype = vval.Type()
		vkind = vtype.Kind()
		ptr = true
	} else {
		vtype = reflect.TypeOf(v)
	}
	
	switch vkind {
	case reflect.Array, reflect.Slice:
		if vval.IsNil() {
			s.writeNil()
		} else {
			s.writePtr(ptr)
		}
	case reflect.Map:
		if vval.IsNil() {
			s.writeNil()
		} else {
			s.writePtr(ptr)
			s.writeMap(vval, vtype)
		}
	case reflect.Struct:
		s.writePtr(ptr)
		s.writeStruct(vval, vtype)
	case reflect.String:
		s.varbuff.WriteString(fmt.Sprintf("\"%v\"", vval.Interface()))
	default:
		s.varbuff.WriteString(fmt.Sprintf("%v", vval.Interface()))
	}
}

var test interface{}

func Marshal(v map[string]interface{}, pkgname string) ([]byte, error) {
	fmt.Println("Serialising.")

	s := &serialiser{
		buffer: bytes.NewBuffer(make([]byte, 0, 256)),
		varbuff: bytes.NewBuffer(make([]byte, 0, 256)),
		indentStr: "\t",
	}

	s.init(pkgname)

	for varname, val := range v {
		s.startVar(varname)
		s.startSerialise(val)
	}

	return s.buffer.Bytes(), nil
}
