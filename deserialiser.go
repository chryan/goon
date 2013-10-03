package goon

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

type TypeFactory interface {
	New(typename, pkgname string) interface{}
}

type Errors struct {
	Msgs []string
}

type deserialiser struct {
	fileset *token.FileSet
	typefactory TypeFactory
	errors []string
}

// Used to keep track of array element token positions.
type arrelement struct {
	item interface{}
	pos token.Pos
}

type mapelement struct {
	key interface{}
	kpos token.Pos
	val interface{}
	vpos token.Pos
}

func (d *deserialiser) deserialiseStruct(typename, pkgname string, elts []ast.Expr) interface{} {
	newstruct := d.typefactory.New(typename, pkgname)

	if newstruct == nil {
		return nil
	}

	rval := reflect.ValueOf(newstruct).Elem() // It's always a pointer to our data.
	rtype := rval.Type()
	kind := rtype.Kind()

	// Make sure we're actually deserialising a struct.
	if kind != reflect.Struct || rtype.Name() != typename {
		return nil
	}

	for _, expr := range elts {
		// Make sure it's a key value pair as intended for struct elements.
		kvexpr, ok := expr.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		ident, ok := kvexpr.Key.(*ast.Ident)
		if !ok {
			continue
		}

		// Look for our field
		fval := rval.FieldByName(ident.Name)

		if !fval.IsValid() || !fval.CanSet() || !ok {
			continue
		}

		if val, pos := d.deserialise(kvexpr.Value); val != nil {
			if ftype, ok := rtype.FieldByName(ident.Name); ok {
				d.assignValue(val, fval, ftype.Type, pos)
			}
		}
	}

	return newstruct
}

// Handle array/slice types.
func (d *deserialiser) deserialiseArray(elts []ast.Expr) []interface{} {
	if len(elts) == 0 {
		return nil
	}

	arrvals := make([]interface{}, 0, len(elts))
	for _, elt := range elts {
		if val, pos := d.deserialise(elt); val != nil {
			arrvals = append(arrvals, &arrelement{val, pos})
		}
	}
	return arrvals
}

func (d *deserialiser) deserialiseMap(elts []ast.Expr) []interface{} {
	if len(elts) == 0 {
		return nil
	}

	mappairs := make([]interface{}, 0, len(elts))
	for _, elt := range elts {
		kvexpr, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		var key interface{}
		var kpos token.Pos

		// Make sure we're dealing with 
		if keyident, ok := kvexpr.Key.(*ast.Ident); ok {
			key, kpos = keyident.Name, keyident.NamePos
		} else {
			key, kpos = d.deserialise(kvexpr.Key)
		}

		if key != nil {
			if val, vpos := d.deserialise(kvexpr.Value); val != nil {
				mappairs = append(mappairs, &mapelement{key, kpos, val, vpos})
			}
		}
	}
	return mappairs
}

// Composites consist of arrays and structs.
func (d *deserialiser) deserialiseComposite(c *ast.CompositeLit) interface{} {
	switch t := c.Type.(type) {
	// Standard type.
	case *ast.Ident:
		if d.typefactory != nil {
			return d.deserialiseStruct(t.Name, "", c.Elts)
		} else {
			return d.deserialiseMap(c.Elts)
		}
	// Type in package.
	case *ast.SelectorExpr:
		if pkg, ok := t.X.(*ast.Ident); ok {
			if d.typefactory != nil {
				return d.deserialiseStruct(t.Sel.Name, pkg.Name, c.Elts)
			} else {
				return d.deserialiseMap(c.Elts)
			}
		}
	case *ast.ArrayType:
		return d.deserialiseArray(c.Elts)
	case *ast.MapType:
		return d.deserialiseMap(c.Elts)
	}

	return nil
}

// Used to set values to a reflection value object.
func (d *deserialiser) assignValue(inval interface{}, outputval reflect.Value, outtype reflect.Type, pos token.Pos) (success bool) {
	inputval := reflect.ValueOf(inval)

	fkind := outputval.Kind()
	setkind := inputval.Kind()

	// Instantiate our type if it's a pointer.
	if fkind == reflect.Ptr && fkind != setkind && outtype != nil {
		newptr := reflect.New(outtype.Elem())
		outputval.Set(newptr)
		outputval = newptr.Elem()
		fkind = outputval.Kind()
	}

	// Recover from type setting failures.
	defer func() {
		if r := recover(); r != nil {
			success = false
			d.errors = append(d.errors,
				fmt.Sprintf("goon: Unable to assign %v value to %v (%v)",
					setkind, outputval.Type().Name(), d.fileset.Position(pos)))
		}
	}()

	// Try to assign the standard types.
	switch fkind {
		case reflect.Float32, reflect.Float64:
			outputval.SetFloat(inputval.Float())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			outputval.SetUint(inputval.Uint())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			outputval.SetInt(inputval.Int())
		case reflect.Bool:
			outputval.SetBool(inputval.Bool())
		case reflect.Struct:
			outputval.Set(inputval.Elem())
		case reflect.Slice:
			slicelen := inputval.Len()
			newslice := reflect.MakeSlice(outtype, 0, slicelen)
			for i := 0; i < slicelen; i++ {
				if _inval, ok := inputval.Index(i).Interface().(*arrelement); ok {
					_outval := reflect.New(outtype.Elem()).Elem()
					if d.assignValue(_inval.item, _outval, _outval.Type(), _inval.pos) && _outval.IsValid() {
						newslice = reflect.Append(newslice, _outval)
					}
				}
			}
			outputval.Set(newslice)
		case reflect.Map:
			slicelen := inputval.Len()
			newmap := reflect.MakeMap(outtype)
			for i := 0; i < slicelen; i++ {
				if _inval, ok := inputval.Index(i).Interface().(*mapelement); ok {
					_outkey := reflect.New(outtype.Key()).Elem()
					_outval := reflect.New(outtype.Elem()).Elem()
					if d.assignValue(_inval.key, _outkey, _outkey.Type(), _inval.kpos) && _outkey.IsValid() {
						if d.assignValue(_inval.val, _outval, _outval.Type(), _inval.kpos) && _outval.IsValid() {
							newmap.SetMapIndex(_outkey, _outval)
						}
					}
				}
			}
			outputval.Set(newmap)
		default:
			outputval.Set(inputval)
	}

	success = true
	return
}

// Main recursive deserialise function.
// Handles string, float, int types and recurses into complex types if needed.
func (d *deserialiser) deserialise(astval interface{}) (retval interface{}, pos token.Pos) {
	if actual, ok := astval.(*ast.UnaryExpr); ok {
		// Reassign to the actual type.
		astval = actual.X
	}

	switch t := astval.(type) {
	case *ast.CompositeLit:
		pos = t.Lbrace
		retval = d.deserialiseComposite(t)
	case *ast.BasicLit:
		pos	= t.ValuePos
		switch t.Kind {
		case token.STRING:
			retval = strings.Trim(t.Value, "\"")
		case token.INT:
			if i, err := strconv.Atoi(t.Value); err == nil {
				retval = i
			}
		case token.FLOAT:
			if f, err := strconv.ParseFloat(t.Value, 64); err == nil {
				retval = f
			}
		}
	// Booleans are treated this way.
	case *ast.Ident:
		if b, err := strconv.ParseBool(t.Name); err == nil {
			retval = b
		}
	}

	return
}

func Unmarshal(filename string, data []byte) (map[string]interface{}, *Errors) {
	return UnmarshalTyped(filename, data, nil)
}

func UnmarshalTyped(filename string, data []byte, tf TypeFactory) (deserialised map[string]interface{}, errs *Errors) {
	ds := &deserialiser{
		fileset: token.NewFileSet(),
		typefactory: tf,
		errors: make([]string, 0, 8),
	}

	var parsedata interface{}
	if len(data) > 0 {
		parsedata = data
	}

	f, err := parser.ParseFile(ds.fileset, filename, parsedata, 0)
	if err != nil {
		errs = &Errors{[]string{fmt.Sprintf("%v", err)}}
		return
	}

	deserialised = make(map[string]interface{})

	for _, decl := range f.Decls {

		d, ok := decl.(*ast.GenDecl)
		if !ok || d.Tok != token.VAR || len(d.Specs) == 0 {
			continue
		}

		valspec, ok := d.Specs[0].(*ast.ValueSpec)
		if !ok {
			continue
		}

		name := valspec.Names[0].Name
		val, _ := ds.deserialise(valspec.Values[0])
		deserialised[name] = val
	}

	if len(ds.errors) > 0 {
		errs = &Errors{ds.errors}
	}

	return
}