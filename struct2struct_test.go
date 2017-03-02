package struct2struct

import (
	"errors"
	"reflect"
	"testing"
)

type Untagged struct {
	MatchString          string
	MappedNameString     string
	MappedShortPkgString string
	MappedPkgPathString  string
}

type TwoIntsA struct {
	First  int
	Second int `TwoIntsB:"SecondB"`
}

type TwoIntsB struct {
	SecondB int
	First   int
}

func TestMarshal(t *testing.T) {
	var tests = []struct {
		name     string
		in       interface{}
		other    interface{}
		expected interface{}
		err      error
	}{
		{
			name: "Untagged, other nil",
			in: struct {
				Str string
			}{
				Str: "string",
			},
			err: errors.New("nil target"),
		},
		{
			name: "Tagged - partial",
			in: struct {
				MatchString string
			}{
				MatchString: "match",
			},
			other: &Untagged{},
			expected: &Untagged{
				MatchString: "match",
			},
		},
		{
			name: "Tagged - complete",
			in: struct {
				MatchString    string
				NameString     string `Untagged:"MappedNameString"`
				ShortPkgString string `struct2struct.Untagged:"MappedShortPkgString"`
				PkgPathString  string `github.com/theothertomelliott/struct2struct.Untagged:"MappedPkgPathString"`
			}{
				MatchString:    "match",
				NameString:     "name",
				ShortPkgString: "shortPkg",
				PkgPathString:  "pkgPath",
			},
			other: &Untagged{},
			expected: &Untagged{
				MatchString:          "match",
				MappedNameString:     "name",
				MappedShortPkgString: "shortPkg",
				MappedPkgPathString:  "pkgPath",
			},
		},
		{
			name: "Target not a pointer",
			in: struct {
				MatchString    string
				NameString     string `Untagged:"MappedNameString"`
				ShortPkgString string `struct2struct.Untagged:"MappedShortPkgString"`
				PkgPathString  string `github.com/theothertomelliott/struct2struct.Untagged:"MappedPkgPathString"`
			}{
				MatchString:    "match",
				NameString:     "name",
				ShortPkgString: "shortPkg",
				PkgPathString:  "pkgPath",
			},
			other: Untagged{},
			err:   errors.New("expect target to be a pointer"),
		},
		{
			name: "Non-matching types",
			in: struct {
				MatchString int
			}{
				MatchString: 100,
			},
			other: &Untagged{},
			err:   errors.New("MatchString: types do not match"),
		},
		{
			name: "Unsettable field",
			in: struct {
				MatchString int
				unexported  int
			}{
				MatchString: 100,
				unexported:  200,
			},
			other: &struct {
				MatchString int
				unexported  int
			}{},
			expected: &struct {
				MatchString int
				unexported  int
			}{
				MatchString: 100,
			},
		},
		{
			name: "String pointer to string pointer",
			in: struct {
				MatchString *string
			}{
				MatchString: stringPtr("match"),
			},
			other: &struct {
				MatchString *string
			}{},
			expected: &struct {
				MatchString *string
			}{
				MatchString: stringPtr("match"),
			},
		},
		{
			name: "String pointer to string",
			in: struct {
				MatchString *string
			}{
				MatchString: stringPtr("match"),
			},
			other: &struct {
				MatchString string
			}{},
			expected: &struct {
				MatchString string
			}{
				MatchString: "match",
			},
		},
		{
			name: "String to string pointer",
			in: struct {
				MatchString string
			}{
				MatchString: "match",
			},
			other: &struct {
				MatchString *string
			}{},
			expected: &struct {
				MatchString *string
			}{
				MatchString: stringPtr("match"),
			},
		},
		{
			name: "String to int pointer",
			in: struct {
				MatchString string
			}{
				MatchString: "match",
			},
			other: &struct {
				MatchString *int
			}{},
			err: errors.New("MatchString: types do not match"),
		},
		{
			name: "Struct field, matching",
			in: struct {
				SubStruct struct{ num int }
			}{
				SubStruct: struct{ num int }{num: 100},
			},
			other: &struct {
				SubStruct struct{ num int }
			}{},
			expected: &struct {
				SubStruct struct{ num int }
			}{
				SubStruct: struct{ num int }{num: 100},
			},
		},
		{
			name: "Struct field, not matching",
			in: struct {
				SubStruct TwoIntsA
			}{
				SubStruct: TwoIntsA{
					First:  10,
					Second: 20,
				},
			},
			other: &struct {
				SubStruct TwoIntsB
			}{},
			expected: &struct {
				SubStruct TwoIntsB
			}{
				SubStruct: TwoIntsB{
					SecondB: 20,
					First:   10,
				},
			},
		},
		{
			name: "Struct field, pointer to pointer",
			in: struct {
				SubStruct *TwoIntsA
			}{
				SubStruct: &TwoIntsA{
					First:  10,
					Second: 20,
				},
			},
			other: &struct {
				SubStruct *TwoIntsB
			}{},
			expected: &struct {
				SubStruct *TwoIntsB
			}{
				SubStruct: &TwoIntsB{
					SecondB: 20,
					First:   10,
				},
			},
		},
		{
			name: "Struct field, pointer to non-pointer",
			in: struct {
				SubStruct *TwoIntsA
			}{
				SubStruct: &TwoIntsA{
					First:  10,
					Second: 20,
				},
			},
			other: &struct {
				SubStruct TwoIntsB
			}{},
			expected: &struct {
				SubStruct TwoIntsB
			}{
				SubStruct: TwoIntsB{
					SecondB: 20,
					First:   10,
				},
			},
		},
		{
			name: "Struct field, non-pointer to pointer",
			in: struct {
				SubStruct TwoIntsA
			}{
				SubStruct: TwoIntsA{
					First:  10,
					Second: 20,
				},
			},
			other: &struct {
				SubStruct *TwoIntsB
			}{},
			expected: &struct {
				SubStruct *TwoIntsB
			}{
				SubStruct: &TwoIntsB{
					SecondB: 20,
					First:   10,
				},
			},
		},
		{
			name: "Struct fields, error",
			in: struct {
				SubStruct struct{ First string }
			}{
				SubStruct: struct{ First string }{First: "first"},
			},
			other: &struct {
				SubStruct TwoIntsB
			}{},
			err: errors.New("SubStruct: First: types do not match"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Marshal(
				test.in,
				test.other,
			)
			if test.err == nil && err != nil {
				t.Error(err)
			}
			if test.err != nil && err == nil {
				t.Error("expected an error")
			}
			if test.err != nil && err != nil && test.err.Error() != err.Error() {
				t.Errorf("errors did not match, expected '%v', got '%v'", test.err, err)
			}
			if err == nil && !reflect.DeepEqual(test.expected, test.other) {
				t.Errorf("values did not match, expected '%v', got '%v'", test.expected, test.other)
			}
		})
	}
}

func TestMarshalStrict(t *testing.T) {
	err := MarshalStrict(nil, nil)
	if err == nil || err.Error() != "not implemented" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMapFields(t *testing.T) {
	var tests = []struct {
		name           string
		in             interface{}
		other          interface{}
		expectedValues map[string]string
	}{
		{
			name: "Untagged, other nil",
			in: struct {
				Str string
			}{
				Str: "string",
			},
			other: nil,
			expectedValues: map[string]string{
				"Str": "string",
			},
		},
		{
			name: "Tagged",
			in: struct {
				MatchString    string
				NameString     string `Untagged:"MappedNameString"`
				ShortPkgString string `struct2struct.Untagged:"MappedShortPkgString"`
				PkgPathString  string `github.com/theothertomelliott/struct2struct.Untagged:"MappedPkgPathString"`
			}{
				MatchString:    "match",
				NameString:     "name",
				ShortPkgString: "shortPkg",
				PkgPathString:  "pkgPath",
			},
			other: Untagged{},
			expectedValues: map[string]string{
				"MatchString":          "match",
				"MappedNameString":     "name",
				"MappedShortPkgString": "shortPkg",
				"MappedPkgPathString":  "pkgPath",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := mapFields(
				test.in,
				test.other,
			)
			if len(test.expectedValues) != len(result) {
				t.Errorf("unexpected number of fields in result: %v", len(result))
			}
			for name, value := range test.expectedValues {
				if v, ok := result[name]; ok {
					if v.String() != value {
						t.Errorf("incorrect value for %v. Expected %v, got %v", name, value, v.String())
					}
				} else {
					t.Errorf("expected value not mapped: %v", name)
				}
			}
		})
	}
}

func stringPtr(in string) *string {
	return &in
}
