package goon

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

type serialiser struct {
	buffer  *bytes.Buffer
	varbuff *bytes.Buffer

	indent    int
	indentStr string
	pkgName   string
}

func (s *serialiser) startVar(name string) {
	s.buffer.WriteString(fmt.Sprintf("var %v ", name))
}

func (s *serialiser) startSerialise(v interface{}) {
	s.serialise(reflect.ValueOf(v))

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
	s.pkgName = pkgname
	s.buffer.WriteString(fmt.Sprintf("package %v\n\n", s.pkgName))
}

func (s *serialiser) incIndent() {
	s.indent++
}

func (s *serialiser) decIndent() {
	s.indent--
}

func (s *serialiser) writeIndents() {
	for i := 0; i < s.indent; i++ {
		s.varbuff.WriteString(s.indentStr)
	}
}

func (s *serialiser) writePtr(ptr bool) {
	if ptr {
		s.varbuff.WriteString("&")
	}
}

func (s *serialiser) writeNil() {
	s.varbuff.WriteString("nil")
}

func (s *serialiser) getTypeName(typ reflect.Type) string {
	if typename := typ.Name(); len(typename) > 0 {
		if fullpkgname := typ.PkgPath(); len(fullpkgname) > 0 {
			paths := strings.Split(fullpkgname, "/")
			pkgname := paths[len(paths)-1]
			if pkgname != s.pkgName {
				return pkgname + "." + typename
			}
		}
		return typename
	}
	return "interface{}"
}

func (s *serialiser) serialiseStruct(vval reflect.Value, vtype reflect.Type) {
	s.varbuff.WriteString(s.getTypeName(vtype))
	s.varbuff.WriteString("{")

	if numfield := vtype.NumField(); numfield > 0 {
		s.incIndent()
		for i := 0; i < numfield; i++ {
			ftype := vtype.Field(i)
			fval := vval.Field(i)

			s.varbuff.WriteString("\n")
			s.writeIndents()
			s.varbuff.WriteString(fmt.Sprintf("%v: ", ftype.Name))
			s.serialise(fval)
			s.varbuff.WriteString(",")
		}
		s.varbuff.WriteString("\n")
		s.decIndent()
		s.writeIndents()
	}
	s.varbuff.WriteString("}")
}

func (s *serialiser) serialiseSeq(vval reflect.Value, vtype reflect.Type) {
	s.varbuff.WriteString("[]")
	s.varbuff.WriteString(s.getTypeName(vtype))
	s.varbuff.WriteString("{")
	if seqlen := vval.Len(); seqlen > 0 {
		s.incIndent()
		for i := 0; i < seqlen; i++ {
			sval := vval.Index(i)
			s.varbuff.WriteString("\n")
			s.writeIndents()
			s.serialise(sval)
			s.varbuff.WriteString(",")
		}
		s.varbuff.WriteString("\n")
		s.decIndent()
		s.writeIndents()
	}
	s.varbuff.WriteString("}")
}

func (s *serialiser) serialiseMap(vval reflect.Value, vtype reflect.Type) {
	s.varbuff.WriteString(fmt.Sprintf("map[%v]%v", s.getTypeName(vtype.Key()), s.getTypeName(vtype.Elem())))
	s.varbuff.WriteString("{")

	keys := vval.MapKeys()

	if len(keys) > 0 {
		s.incIndent()
		for _, key := range keys {
			val := vval.MapIndex(key)
			if key.CanInterface() && val.CanInterface() {
				s.varbuff.WriteString("\n")
				s.writeIndents()
				s.serialise(key)
				s.varbuff.WriteString(": ")
				s.serialise(val)
				s.varbuff.WriteString(",")
			}
		}
		s.varbuff.WriteString("\n")
		s.decIndent()
		s.writeIndents()
	}
	s.varbuff.WriteString("}")
}

func (s *serialiser) serialise(vval reflect.Value) {
	vkind := vval.Kind()
	ptr := false

	if vkind == reflect.Ptr || vkind == reflect.Interface {
		if vval.IsNil() {
			s.writeNil()
			return
		}
		vval = vval.Elem()
		ptr = true
	}

	vtype := vval.Type()
	vkind = vtype.Kind()

	switch vkind {
	case reflect.Array, reflect.Slice:
		if vval.IsNil() {
			s.writeNil()
		} else {
			s.writePtr(ptr)
			s.serialiseSeq(vval, vtype)
		}
	case reflect.Map:
		if vval.IsNil() {
			s.writeNil()
		} else {
			s.writePtr(ptr)
			s.serialiseMap(vval, vtype)
		}
	case reflect.Struct:
		s.writePtr(ptr)
		s.serialiseStruct(vval, vtype)
	case reflect.String:
		s.varbuff.WriteString(fmt.Sprintf("\"%v\"", vval.String()))
	case reflect.Float32, reflect.Float64:
		s.varbuff.WriteString(fmt.Sprintf("%g", vval.Float()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s.varbuff.WriteString(fmt.Sprintf("%v", vval.Uint()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s.varbuff.WriteString(fmt.Sprintf("%v", vval.Int()))
	case reflect.Bool:
		s.varbuff.WriteString(fmt.Sprintf("%v", vval.Bool()))
	default:
		s.serialise(reflect.ValueOf(vval.Interface()))
	}
}

var test interface{}

func Marshal(v map[string]interface{}, pkgname string) ([]byte, error) {
	s := &serialiser{
		buffer:    bytes.NewBuffer(make([]byte, 0, 256)),
		varbuff:   bytes.NewBuffer(make([]byte, 0, 256)),
		indentStr: "\t",
	}

	s.init(pkgname)

	for varname, val := range v {
		s.startVar(varname)
		s.startSerialise(val)
	}

	return s.buffer.Bytes(), nil
}
