package struct2struct

import (
	"errors"
	"fmt"
	"reflect"
)

var appliers []applier

func init() {
	appliers = []applier{
		matchedTypeApplier,
		structApplier,
		pointerApplier,
	}
}

type applier func(reflect.Value, reflect.Value) (bool, error)

func applyField(iField reflect.Value, vField reflect.Value) error {
	if !vField.CanSet() {
		return nil
	}
	for _, applier := range appliers {
		applied, err := applier(iField, vField)
		if applied || err != nil {
			return err
		}
	}
	return errors.New("could not apply types")
}

func matchedTypeApplier(iField reflect.Value, vField reflect.Value) (bool, error) {
	if iField.Type() == vField.Type() {
		vField.Set(iField)
		return true, nil
	}
	return false, nil
}

func structApplier(iField reflect.Value, vField reflect.Value) (bool, error) {
	if !(iField.Type().Kind() == reflect.Struct && vField.Type().Kind() == reflect.Struct) {
		return false, nil
	}
	newPtr := reflect.New(vField.Type())
	newPtr.Elem().Set(vField)
	err := marshalStruct(iField.Interface(), newPtr.Interface())
	vField.Set(newPtr.Elem())
	return err == nil, err
}

func marshalStruct(i interface{}, v interface{}) error {
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

func pointerApplier(iField reflect.Value, vField reflect.Value) (bool, error) {
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
		t := reflect.TypeOf(vField.Interface())
		if iField.Kind() == reflect.Struct && t.Elem().Kind() == reflect.Struct {
			newPtr := reflect.New(t.Elem())
			err := applyField(iField, newPtr.Elem())
			if err == nil {
				vField.Set(newPtr)
			}
			return err == nil, err
		}
	}
	return false, nil
}
