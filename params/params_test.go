package params_test

import (
	"encoding/json"
	"fmt"

	"github.com/0x51-dev/jsonrpc2/params"
)

func ExampleGetMap() {
	var nok = map[string]int{
		"A": 0,
		"B": 1,
		"C": 2,
	}
	var ok = map[string]any{
		"A": 0,
		"B": 1,
		"C": 2,
	}
	fmt.Println(params.GetMap[int](nok))
	fmt.Println(params.GetMap[int](ok))
	// Output:
	// map[] false
	// map[A:0 B:1 C:2] true
}

func ExampleGet() {
	fmt.Println(params.Get[string](nil))
	fmt.Println(params.Get[string]("ok"))
	// Output:
	// false
	// ok true
}

func ExampleNumber() {
	fmt.Println(params.GetInt64(json.Number("1")))
	fmt.Println(params.GetInt64(json.Number("2.")))
	fmt.Println(params.GetFloat64(json.Number("3")))
	fmt.Println(params.GetFloat64(json.Number("4.5")))
	// Output:
	// 1 true
	// 0 false
	// 3 true
	// 4.5 true
}
