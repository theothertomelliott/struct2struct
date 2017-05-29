package struct2struct_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/theothertomelliott/struct2struct"
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

type marshalTest struct {
	name     string
	in       interface{}
	other    interface{}
	expected interface{}
	err      error
}

func TestMarshalStructs(t *testing.T) {
	var tests = []marshalTest{
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
				ShortPkgString string `struct2struct_test.Untagged:"MappedShortPkgString"`
				PkgPathString  string `github.com/theothertomelliott/struct2struct_test.Untagged:"MappedPkgPathString"`
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
			err:   errors.New("MatchString: could not apply types"),
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
			err: errors.New("MatchString: could not apply types"),
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
			err: errors.New("SubStruct: First: could not apply types"),
		},
	}
	executeTests(t, tests)
}

func TestMarshalSlices(t *testing.T) {
	var tests = []marshalTest{
		{
			name: "Matching slice types",
			in: []string{
				"a", "b",
			},
			other: &[]string{},
			expected: &[]string{
				"a", "b",
			},
		},
		{
			name: "Non-matching slice types",
			in: []TwoIntsA{
				{
					First:  10,
					Second: 20,
				},
			},
			other: &[]TwoIntsB{},
			expected: &[]TwoIntsB{
				{
					First:   10,
					SecondB: 20,
				},
			},
		},
		{
			name: "Non-matching slice types in struct",
			in: struct {
				Arr []TwoIntsA
			}{
				Arr: []TwoIntsA{
					{
						First:  10,
						Second: 20,
					},
				},
			},
			other: &struct {
				Arr []TwoIntsB
			}{},
			expected: &struct {
				Arr []TwoIntsB
			}{
				Arr: []TwoIntsB{
					{
						First:   10,
						SecondB: 20,
					},
				},
			},
		},
		{
			name: "Slice to non-slice error",
			in: []string{
				"a", "b",
			},
			other: &struct{}{},
			err:   errors.New("cannot apply a slice to a non-slice value"),
		},
		{
			name: "Non-slice to slice error",
			in:   &struct{}{},
			other: &[]string{
				"a", "b",
			},
			err: errors.New("cannot apply a non-slice value to a slice"),
		},
	}
	executeTests(t, tests)
}

func TestMarshalMaps(t *testing.T) {
	var tests = []marshalTest{
		{
			name: "Matching slice types",
			in: []string{
				"a", "b",
			},
			other: &[]string{},
			expected: &[]string{
				"a", "b",
			},
		},
	}
	executeTests(t, tests)
}

func executeTests(t *testing.T, tests []marshalTest) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := struct2struct.Marshal(
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

func stringPtr(in string) *string {
	return &in
}
