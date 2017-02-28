package struct2struct

import (
	"fmt"
	"reflect"
	"testing"
)

type Example struct {
	Stuff      string `struct2struct.ExampleTo:"OtherStuff"`
	Nonsense   string `github.com/theothertomelliott/struct2struct.ExampleTo:"OtherNonsense"`
	unexported int32
}

type ExampleTo struct {
	OtherStuff    string
	OtherNonsense string
}

func TestMapFields(t *testing.T) {
	result := mapFields(Example{Stuff: "blah", Nonsense: "bloo"}, reflect.TypeOf(ExampleTo{}))
	fmt.Println(result)
}
