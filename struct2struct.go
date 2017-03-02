package struct2struct

import (
	"errors"
	"fmt"
	"reflect"
)

// Marshal processes i and applies its values to v.
// Fields are matched first by s2s tags, then by field names.
func Marshal(i interface{}, v interface{}) error {
	return doMarshal(i, v, false)
}

// MarshalStrict processes i and applies its values to v as with Marhsal.
// If any values in i are not converted, an error will be thrown.
func MarshalStrict(i interface{}, v interface{}) error {
	return errors.New("not implemented")
}

func doMarshal(i interface{}, v interface{}, strict bool) error {
	if v == nil {
		return errors.New("nil target")
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("expect target to be a pointer")
	}

	iFields := mapFields(i, v)
	vFields := mapFields(v, i)

	for name, iField := range iFields {
		if vField, ok := vFields[name]; ok {
			err := applyField(iField, vField)
			if err != nil {
				return fmt.Errorf("%v: %v", name, err)
			}
		}
	}

	return nil
}

func applyField(iField reflect.Value, vField reflect.Value) error {
	if !vField.CanSet() {
		return nil
	}
	if iField.Type() != vField.Type() {
		applied, err := applyStruct(iField, vField)
		if applied || err != nil {
			return err
		}
		applied, err = applyPointers(iField, vField)
		if applied || err != nil {
			return err
		}
		return errors.New("types do not match")
	}

	vField.Set(iField)
	return nil
}

func applyPointers(iField reflect.Value, vField reflect.Value) (bool, error) {
	if iField.Type().Kind() == reflect.Ptr {
		err := applyField(reflect.Indirect(iField), vField)
		return err == nil, err
	}
	iPtrType := reflect.PtrTo(iField.Type())
	if vField.Type().Kind() == reflect.Ptr {
		if iPtrType == vField.Type() {
			newPtr := reflect.New(iField.Type())
			newPtr.Elem().Set(iField)
			err := applyField(newPtr, vField)
			return err == nil, err
		}
	}
	return false, nil
}

func applyStruct(iField reflect.Value, vField reflect.Value) (bool, error) {
	if !(iField.Type().Kind() == reflect.Struct && vField.Type().Kind() == reflect.Struct) {
		return false, nil
	}
	newPtr := reflect.New(vField.Type())
	newPtr.Elem().Set(vField)
	err := Marshal(iField.Interface(), newPtr.Interface())
	vField.Set(newPtr.Elem())
	return err == nil, err
}

func mapFields(i interface{}, other interface{}) map[string]reflect.Value {

	var outFields = make(map[string]reflect.Value)
	iValue := reflect.Indirect(reflect.ValueOf(i))
	iType := iValue.Type()

	var otherType reflect.Type
	if other != nil {
		otherValue := reflect.ValueOf(other)
		if reflect.TypeOf(other).Kind() == reflect.Ptr {
			otherValue = reflect.Indirect(otherValue)
		}
		otherType = otherValue.Type()
	}

	for i := 0; i < iValue.NumField(); i++ {
		fType := iType.Field(i)
		fValue := iValue.Field(i)
		tags := fType.Tag
		if otherType != nil {
			if name, ok := tags.Lookup(fmt.Sprintf("%v.%v", otherType.PkgPath(), otherType.Name())); ok {
				outFields[name] = fValue
				continue
			}
			if name, ok := tags.Lookup(otherType.String()); ok {
				outFields[name] = fValue
				continue
			}
			if name, ok := tags.Lookup(otherType.Name()); ok {
				outFields[name] = fValue
				continue
			}
		}
		outFields[iType.Field(i).Name] = fValue
	}
	return outFields
}

// Custom allows a struct to provide custom marshalling to another struct type.
// Custom marshaling will be performed after automatic marshaling.
type Custom interface {
	MarshalStruct(v interface{}) error
	UnmarshalStruct(v interface{}) error
}
