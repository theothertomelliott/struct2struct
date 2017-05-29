package struct2struct

import "testing"

type Untagged struct {
	MatchString          string
	MappedNameString     string
	MappedShortPkgString string
	MappedPkgPathString  string
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
