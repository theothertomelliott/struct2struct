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
	return doMarshal(i, v, true)
}

func doMarshal(i interface{}, v interface{}, strict bool) error {
	return errors.New("not implemented")
}

func mapFields(i interface{}, otherType reflect.Type) map[string]reflect.Value {
	var outFields = make(map[string]reflect.Value)
	iValue := reflect.ValueOf(i)
	iType := reflect.TypeOf(i)

	for i := 0; i < iValue.NumField(); i++ {
		fType := iType.Field(i)
		fValue := iValue.Field(i)
		tags := fType.Tag
		if name, ok := tags.Lookup(fmt.Sprintf("%v.%v", otherType.PkgPath(), otherType.Name())); ok {
			outFields[name] = fValue
			continue
		}
		if name, ok := tags.Lookup(otherType.String()); ok {
			outFields[name] = fValue
			continue
		}
		outFields[iType.Field(i).Name] = fValue
	}
	return outFields
}

// Custom allows a struct to provide custom marshalling to another struct type.
// Custom marshaling will be performed after automatic marshaling.
type Custom interface {
	Marshal(v interface{}) error
	Unmarshal(v interface{}) error
}
