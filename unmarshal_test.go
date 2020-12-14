package sfv_test

import (
	"fmt"

	"github.com/ucarion/sfv"
)

func ExampleUnmarshal_raw_item() {
	var item sfv.Item
	fmt.Println(sfv.Unmarshal("text/html; charset=UTF-8", &item))
	fmt.Println(item)
	// Output:
	// <nil>
	// {text/html {map[charset:UTF-8] [charset]}}
}

func ExampleUnmarshal_raw_list() {
	var list []sfv.Member
	fmt.Println(sfv.Unmarshal("fr-CH, fr;q=0.9, *;q=0.5", &list))

	for _, m := range list {
		fmt.Println(m.Item)
	}

	// Output:
	// <nil>
	// {fr-CH {map[] []}}
	// {fr {map[q:0.9] [q]}}
	// {* {map[q:0.5] [q]}}
}

func ExampleUnmarshal_raw_dict() {
	var dict sfv.Dictionary
	fmt.Println(sfv.Unmarshal("public, max-age=604800, immutable", &dict))

	for _, k := range dict.Keys {
		fmt.Println(k, dict.Map[k].Item)
	}

	// Output:
	// <nil>
	// public {true {map[] []}}
	// max-age {604800 {map[] []}}
	// immutable {true {map[] []}}
}
