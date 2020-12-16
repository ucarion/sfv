package sfv_test

import (
	"fmt"

	"github.com/ucarion/sfv"
)

func ExampleMarshal_raw_item() {
	item := sfv.Item{
		BareItem: sfv.BareItem{
			Type:  sfv.BareItemTypeToken,
			Token: "text/html",
		},
		Params: sfv.Params{
			Keys: []string{"charset"},
			Map: map[string]sfv.BareItem{
				"charset": sfv.BareItem{
					Type:  sfv.BareItemTypeToken,
					Token: "UTF-8",
				},
			},
		},
	}

	fmt.Println(sfv.Marshal(item))
	// Output: text/html;charset=UTF-8 <nil>
}

func ExampleMarshal_raw_list() {
	list := sfv.List{
		sfv.Member{
			IsItem: true,
			Item: sfv.Item{
				BareItem: sfv.BareItem{
					Type:  sfv.BareItemTypeToken,
					Token: "fr-CH",
				},
			},
		},
		sfv.Member{
			IsItem: true,
			Item: sfv.Item{
				BareItem: sfv.BareItem{
					Type:  sfv.BareItemTypeToken,
					Token: "fr",
				},
				Params: sfv.Params{
					Keys: []string{"q"},
					Map: map[string]sfv.BareItem{
						"q": sfv.BareItem{
							Type:    sfv.BareItemTypeDecimal,
							Decimal: 0.9,
						},
					},
				},
			},
		},
		sfv.Member{
			IsItem: true,
			Item: sfv.Item{
				BareItem: sfv.BareItem{
					Type:  sfv.BareItemTypeToken,
					Token: "*",
				},
				Params: sfv.Params{
					Keys: []string{"q"},
					Map: map[string]sfv.BareItem{
						"q": sfv.BareItem{
							Type:    sfv.BareItemTypeDecimal,
							Decimal: 0.5,
						},
					},
				},
			},
		},
	}

	fmt.Println(sfv.Marshal(list))
	// Output: fr-CH, fr;q=0.9, *;q=0.5 <nil>
}

func ExampleMarshal_raw_dict() {
	dict := sfv.Dictionary{
		Keys: []string{"public", "max-age", "immutable"},
		Map: map[string]sfv.Member{
			"public": sfv.Member{
				IsItem: true,
				Item: sfv.Item{
					BareItem: sfv.BareItem{
						Type:    sfv.BareItemTypeBoolean,
						Boolean: true,
					},
				},
			},
			"max-age": sfv.Member{
				IsItem: true,
				Item: sfv.Item{
					BareItem: sfv.BareItem{
						Type:    sfv.BareItemTypeInteger,
						Integer: 604800,
					},
				},
			},
			"immutable": sfv.Member{
				IsItem: true,
				Item: sfv.Item{
					BareItem: sfv.BareItem{
						Type:    sfv.BareItemTypeBoolean,
						Boolean: true,
					},
				},
			},
		},
	}

	fmt.Println(sfv.Marshal(dict))
	// Output: public, max-age=604800, immutable <nil>
}

func ExampleMarshal_custom_item() {
	type contentType struct {
		MediaType string
		Charset   string `sfv:"charset"`
		Boundary  string `sfv:"boundary"`
	}

	fmt.Println(sfv.Marshal(contentType{MediaType: "text/html", Charset: "UTF-8"}))
	fmt.Println(sfv.Marshal(contentType{
		MediaType: "multipart/form-data",
		Charset:   "UTF-8",
		Boundary:  "xxx",
	}))

	// Output:
	// text/html;charset=UTF-8 <nil>
	// multipart/form-data;charset=UTF-8;boundary=xxx <nil>
}

func ExampleMarshal_bare_item() {
	fmt.Println(sfv.Marshal("foo"))
	// Output: foo <nil>
}

func ExampleMarshal_custom_basic_list() {
	fmt.Println(sfv.Marshal([]string{"foo", "bar", "baz"}))
	// Output: foo, bar, baz <nil>
}

func ExampleMarshal_custom_list() {
	type language struct {
		Tag    string
		Weight float64 `sfv:"q"`
	}

	fmt.Println(sfv.Marshal([]language{
		language{Tag: "fr-CH"},
		language{Tag: "fr", Weight: 0.9},
		language{Tag: "*", Weight: 0.5},
	}))

	// Output:
	// fr-CH, fr;q=0.9, *;q=0.5 <nil>
}

func ExampleMarshal_custom_list_with_inner_list() {
	fmt.Println(sfv.Marshal([][]string{
		[]string{"gzip", "fr"},
		[]string{"identity", "fr"},
	}))

	// Output:
	// (gzip fr), (identity fr) <nil>
}

func ExampleMarshal_custom_list_with_inner_list_with_params() {
	type innerListWithParams struct {
		Names []string
		Foo   string `sfv:"foo"`
	}

	fmt.Println(sfv.Marshal([]innerListWithParams{
		innerListWithParams{
			Names: []string{"gzip", "fr"},
			Foo:   "bar",
		},
		innerListWithParams{
			Names: []string{"identity", "fr"},
			Foo:   "baz",
		},
	}))

	// Output:
	// (gzip fr);foo=bar, (identity fr);foo=baz <nil>
}

func ExampleMarshal_custom_list_with_inner_list_with_nested_params() {
	type itemWithParams struct {
		Name string
		XXX  string `sfv:"xxx"`
	}

	type innerListWithParams struct {
		Names []itemWithParams
		Foo   string `sfv:"foo"`
	}

	fmt.Println(sfv.Marshal([]innerListWithParams{
		innerListWithParams{
			Names: []itemWithParams{
				itemWithParams{Name: "gzip", XXX: "yyy"},
				itemWithParams{Name: "fr"},
			},
			Foo: "bar",
		},
		innerListWithParams{
			Names: []itemWithParams{
				itemWithParams{Name: "identity"},
				itemWithParams{Name: "fr", XXX: "zzz"},
			},
			Foo: "baz",
		},
	}))

	// Output:
	// (gzip;xxx=yyy fr);foo=bar, (identity fr;xxx=zzz);foo=baz <nil>
}

func ExampleMarshal_custom_basic_map() {
	fmt.Println(sfv.Marshal(map[string]int{"foo": 1}))
	// Output: foo=1 <nil>
}

func ExampleMarshal_custom_map() {
	type thing struct {
		Value int
		Foo   string `sfv:"foo"`
	}

	fmt.Println(sfv.Marshal(map[string]thing{
		"xxx": thing{Value: 1, Foo: "bar"},
	}))

	// Output:
	// xxx=1;foo=bar <nil>
}

func ExampleMarshal_custom_map_inner_list() {
	fmt.Println(sfv.Marshal(map[string][]string{
		"accept-encoding": []string{"gzip", "br"},
	}))

	// Output:
	// accept-encoding=(gzip br) <nil>
}

func ExampleMarshal_custom_map_with_inner_list_with_params() {
	type innerListWithParams struct {
		Names []string
		Foo   string `sfv:"foo"`
	}

	fmt.Println(sfv.Marshal(map[string]innerListWithParams{
		"a": innerListWithParams{
			Names: []string{"gzip", "fr"},
			Foo:   "bar",
		},
	}))

	// Output:
	// a=(gzip fr);foo=bar <nil>
}

func ExampleMarshal_custom_map_with_inner_list_with_nested_params() {
	type itemWithParams struct {
		Name string
		XXX  string `sfv:"xxx"`
	}

	type innerListWithParams struct {
		Names []itemWithParams
		Foo   string `sfv:"foo"`
	}

	fmt.Println(sfv.Marshal(map[string]innerListWithParams{
		"a": innerListWithParams{
			Names: []itemWithParams{
				itemWithParams{Name: "gzip", XXX: "yyy"},
				itemWithParams{Name: "fr"},
			},
			Foo: "bar",
		},
	}))

	// Output:
	// a=(gzip;xxx=yyy fr);foo=bar <nil>
}
