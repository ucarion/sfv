package sfv_test

import (
	"fmt"
	"reflect"
	"testing"

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

func TestMarshal_custom_bare_types(t *testing.T) {
	testCases := []struct {
		Out     string
		Type    reflect.Kind
		Bool    bool
		String  string
		Int     int
		Int8    int8
		Int16   int16
		Int32   int32
		Int64   int64
		Uint    uint
		Uint8   uint8
		Uint16  uint16
		Uint32  uint32
		Uint64  uint64
		Float32 float32
		Float64 float64
		Bytes   []byte
	}{
		{Out: "?0", Type: reflect.Bool, Bool: false},
		{Out: "?1", Type: reflect.Bool, Bool: true},
		{Out: "foo", Type: reflect.String, String: "foo"},
		{Out: "3", Type: reflect.Int, Int: 3},
		{Out: "3", Type: reflect.Int8, Int8: 3},
		{Out: "3", Type: reflect.Int16, Int16: 3},
		{Out: "3", Type: reflect.Int32, Int32: 3},
		{Out: "3", Type: reflect.Int64, Int64: 3},
		{Out: "3", Type: reflect.Uint, Uint: 3},
		{Out: "3", Type: reflect.Uint8, Uint8: 3},
		{Out: "3", Type: reflect.Uint16, Uint16: 3},
		{Out: "3", Type: reflect.Uint32, Uint32: 3},
		{Out: "3", Type: reflect.Uint64, Uint64: 3},
		{Out: "3.14", Type: reflect.Float32, Float32: 3.14},
		{Out: "3.14", Type: reflect.Float64, Float64: 3.14},
		{Out: ":aGVsbG8=:", Type: reflect.Array, Bytes: []byte{'h', 'e', 'l', 'l', 'o'}},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("%s %s", tt.Type, tt.Out), func(t *testing.T) {
			var got interface{}
			var err error

			switch tt.Type {
			case reflect.Bool:
				got, err = sfv.Marshal(tt.Bool)
			case reflect.String:
				got, err = sfv.Marshal(tt.String)
			case reflect.Int:
				got, err = sfv.Marshal(tt.Int)
			case reflect.Int8:
				got, err = sfv.Marshal(tt.Int8)
			case reflect.Int16:
				got, err = sfv.Marshal(tt.Int16)
			case reflect.Int32:
				got, err = sfv.Marshal(tt.Int32)
			case reflect.Int64:
				got, err = sfv.Marshal(tt.Int64)
			case reflect.Uint:
				got, err = sfv.Marshal(tt.Uint)
			case reflect.Uint8:
				got, err = sfv.Marshal(tt.Uint8)
			case reflect.Uint16:
				got, err = sfv.Marshal(tt.Uint16)
			case reflect.Uint32:
				got, err = sfv.Marshal(tt.Uint32)
			case reflect.Uint64:
				got, err = sfv.Marshal(tt.Uint64)
			case reflect.Float32:
				got, err = sfv.Marshal(tt.Float32)
			case reflect.Float64:
				got, err = sfv.Marshal(tt.Float64)
			case reflect.Array:
				got, err = sfv.Marshal(tt.Bytes)
			default:
				panic("bad tt.Type")
			}

			if err != nil {
				t.Errorf("err: %v", err)
				return
			}

			if tt.Out != got {
				t.Errorf("bad unmarshal result: want: %#v, got: %#v", tt.Out, got)
			}
		})
	}
}
