package sfv_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ucarion/sfv"
)

func ExampleUnmarshal_raw_item() {
	var item sfv.Item
	fmt.Println(sfv.Unmarshal("text/html; charset=UTF-8", &item))
	fmt.Println(item.BareItem.Token)
	fmt.Println(item.Params.Map["charset"].Token)

	// Output:
	// <nil>
	// text/html
	// UTF-8
}

func ExampleUnmarshal_raw_list() {
	var list []sfv.Member
	fmt.Println(sfv.Unmarshal("fr-CH, fr;q=0.9, *;q=0.5", &list))

	for _, m := range list {
		fmt.Println(m.Item.BareItem.Token)

		if _, ok := m.Item.Params.Map["q"]; ok {
			fmt.Println(m.Item.Params.Map["q"].Decimal)
		}
	}

	// Output:
	// <nil>
	// fr-CH
	// fr
	// 0.9
	// *
	// 0.5
}

func ExampleUnmarshal_raw_dict() {
	var dict sfv.Dictionary
	fmt.Println(sfv.Unmarshal("public, max-age=604800, immutable", &dict))

	for _, k := range dict.Keys {
		fmt.Println(k, dict.Map[k].Item.BareItem)
	}

	// Output:
	// <nil>
	// public {boolean 0 0   [] true}
	// max-age {integer 604800 0   [] false}
	// immutable {boolean 0 0   [] true}
}

func ExampleUnmarshal_custom_bare_item() {
	var data string
	fmt.Println(sfv.Unmarshal("text/html; charset=UTF-8", &data))
	fmt.Println(data)

	// Output:
	// <nil>
	// text/html
}

func ExampleUnmarshal_custom_item() {
	type contentType struct {
		MediaType string
		Charset   string `sfv:"charset"`
		Boundary  string `sfv:"boundary"`
	}

	var data1 contentType
	fmt.Println(sfv.Unmarshal("text/html; charset=UTF-8", &data1))
	fmt.Println(data1)

	var data2 contentType
	fmt.Println(sfv.Unmarshal("multipart/form-data; charset=UTF-8; boundary=xxx", &data2))
	fmt.Println(data2)

	// Output:
	// <nil>
	// {text/html UTF-8 }
	// <nil>
	// {multipart/form-data UTF-8 xxx}
}

func ExampleUnmarshal_custom_basic_list() {
	var data []string
	fmt.Println(sfv.Unmarshal("foo, bar, baz", &data))
	fmt.Println(data)

	// Output:
	// <nil>
	// [foo bar baz]
}

func ExampleUnmarshal_custom_list() {
	type language struct {
		Tag    string
		Weight float64 `sfv:"q"`
	}

	var data []language
	fmt.Println(sfv.Unmarshal("fr-CH, fr;q=0.9, *;q=0.5", &data))
	fmt.Println(data)

	// Output:
	// <nil>
	// [{fr-CH 0} {fr 0.9} {* 0.5}]
}

func ExampleUnmarshal_custom_list_iterated_calls() {
	type language struct {
		Tag    string
		Weight float64 `sfv:"q"`
	}

	var data []language
	fmt.Println(sfv.Unmarshal("fr-CH", &data))
	fmt.Println(sfv.Unmarshal("fr;q=0.9", &data))
	fmt.Println(sfv.Unmarshal("*;q=0.5", &data))
	fmt.Println(data)

	// Output:
	// <nil>
	// <nil>
	// <nil>
	// [{fr-CH 0} {fr 0.9} {* 0.5}]
}

func ExampleUnmarshal_custom_list_with_inner_list() {
	var data [][]string
	fmt.Println(sfv.Unmarshal("(gzip fr), (identity fr)", &data))
	fmt.Println(data)

	// Output:
	// <nil>
	// [[gzip fr] [identity fr]]
}

func ExampleUnmarshal_custom_list_with_inner_list_with_params() {
	type innerListWithParams struct {
		Names []string
		Foo   string `sfv:"foo"`
	}

	var data []innerListWithParams
	fmt.Println(sfv.Unmarshal("(gzip fr);foo=bar, (identity fr);foo=baz", &data))
	fmt.Println(data)

	// Output:
	// <nil>
	// [{[gzip fr] bar} {[identity fr] baz}]
}

func ExampleUnmarshal_custom_list_with_inner_list_with_nested_params() {
	type itemWithParams struct {
		Name string
		XXX  string `sfv:"xxx"`
	}

	type innerListWithParams struct {
		Names []itemWithParams
		Foo   string `sfv:"foo"`
	}

	var data []innerListWithParams
	fmt.Println(sfv.Unmarshal("(gzip;xxx=yyy fr);foo=bar, (identity fr;xxx=zzz);foo=baz", &data))
	fmt.Println(data)

	// Output:
	// <nil>
	// [{[{gzip yyy} {fr }] bar} {[{identity } {fr zzz}] baz}]
}

func ExampleUnmarshal_custom_basic_map() {
	var data map[string]int
	fmt.Println(sfv.Unmarshal("foo=1, bar=2, baz=3", &data))
	fmt.Println(data)

	// Output:
	// <nil>
	// map[bar:2 baz:3 foo:1]
}

func ExampleUnmarshal_custom_map() {
	type thing struct {
		Value int
		Foo   string `sfv:"foo"`
	}

	var data map[string]thing
	fmt.Println(sfv.Unmarshal("xxx=1;foo=bar, yyy=2;foo=baz", &data))
	fmt.Println(data)

	// Output:
	// <nil>
	// map[xxx:{1 bar} yyy:{2 baz}]
}

func ExampleUnmarshal_custom_map_inner_list() {
	var data map[string][]string
	fmt.Println(sfv.Unmarshal("accept-encoding=(gzip br), accept-language=(en fr)", &data))
	fmt.Println(data)

	// Output:
	// <nil>
	// map[accept-encoding:[gzip br] accept-language:[en fr]]
}

func ExampleUnmarshal_custom_map_with_inner_list_with_params() {
	type innerListWithParams struct {
		Names []string
		Foo   string `sfv:"foo"`
	}

	var data map[string]innerListWithParams
	fmt.Println(sfv.Unmarshal("a=(gzip fr);foo=bar, b=(identity fr);foo=baz", &data))
	fmt.Println(data)

	// Output:
	// <nil>
	// map[a:{[gzip fr] bar} b:{[identity fr] baz}]
}

func ExampleUnmarshal_custom_map_with_inner_list_with_nested_params() {
	type itemWithParams struct {
		Name string
		XXX  string `sfv:"xxx"`
	}

	type innerListWithParams struct {
		Names []itemWithParams
		Foo   string `sfv:"foo"`
	}

	var data map[string]innerListWithParams
	fmt.Println(sfv.Unmarshal("a=(gzip;xxx=yyy fr);foo=bar, b=(identity fr;xxx=zzz);foo=baz", &data))
	fmt.Println(data)

	// Output:
	// <nil>
	// map[a:{[{gzip yyy} {fr }] bar} b:{[{identity } {fr zzz}] baz}]
}

func TestUnmarshal_custom_bare_types(t *testing.T) {
	testCases := []struct {
		In      string
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
		{In: "?0", Type: reflect.Bool, Bool: false},
		{In: "?1", Type: reflect.Bool, Bool: true},
		{In: "?1;foo=bar", Type: reflect.Bool, Bool: true},
		{In: "foo", Type: reflect.String, String: "foo"},
		{In: "\"foo\"", Type: reflect.String, String: "foo"},
		{In: "3", Type: reflect.Int, Int: 3},
		{In: "3", Type: reflect.Int8, Int8: 3},
		{In: "3", Type: reflect.Int16, Int16: 3},
		{In: "3", Type: reflect.Int32, Int32: 3},
		{In: "3", Type: reflect.Int64, Int64: 3},
		{In: "3", Type: reflect.Uint, Uint: 3},
		{In: "3", Type: reflect.Uint8, Uint8: 3},
		{In: "3", Type: reflect.Uint16, Uint16: 3},
		{In: "3", Type: reflect.Uint32, Uint32: 3},
		{In: "3", Type: reflect.Uint64, Uint64: 3},
		{In: "3.14", Type: reflect.Float32, Float32: 3.14},
		{In: "3.14", Type: reflect.Float64, Float64: 3.14},
		{In: ":aGVsbG8=:", Type: reflect.Array, Bytes: []byte{'h', 'e', 'l', 'l', 'o'}},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("%s %s", tt.Type, tt.In), func(t *testing.T) {
			var want interface{}
			var got interface{}
			var err error

			switch tt.Type {
			case reflect.Bool:
				want = tt.Bool
				var v bool
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			case reflect.String:
				want = tt.String
				var v string
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			case reflect.Int:
				want = tt.Int
				var v int
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			case reflect.Int8:
				want = tt.Int8
				var v int8
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			case reflect.Int16:
				want = tt.Int16
				var v int16
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			case reflect.Int32:
				want = tt.Int32
				var v int32
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			case reflect.Int64:
				want = tt.Int64
				var v int64
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			case reflect.Uint:
				want = tt.Uint
				var v uint
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			case reflect.Uint8:
				want = tt.Uint8
				var v uint8
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			case reflect.Uint16:
				want = tt.Uint16
				var v uint16
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			case reflect.Uint32:
				want = tt.Uint32
				var v uint32
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			case reflect.Uint64:
				want = tt.Uint64
				var v uint64
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			case reflect.Float32:
				want = tt.Float32
				var v float32
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			case reflect.Float64:
				want = tt.Float64
				var v float64
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			case reflect.Array:
				want = tt.Bytes
				var v []byte
				err = sfv.Unmarshal(tt.In, &v)
				got = v
			default:
				panic("bad tt.Type")
			}

			if err != nil {
				t.Errorf("err: %v", err)
				return
			}

			if !reflect.DeepEqual(want, got) {
				t.Errorf("bad unmarshal result: want: %#v, got: %#v", want, got)
			}
		})
	}
}

func TestUnmarshal_compound_custom_types(t *testing.T) {
	testCases := []struct {
		input    string
		receiver func() interface{}
		want     interface{}
	}{
		{
			input:    "sig1=:AQIDBA==:, sig2=:AQIDBA==:",
			receiver: func() interface{} { return &map[string][]byte{} },
			want: &map[string][]byte{
				"sig1": {0x01, 0x02, 0x03, 0x04},
				"sig2": {0x01, 0x02, 0x03, 0x04},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(fmt.Sprintf("%T", tt.input), func(t *testing.T) {
			var got interface{}
			var err error

			got = tt.receiver()
			err = sfv.Unmarshal(tt.input, got)

			if err != nil {
				t.Errorf("err: %v", err)
				return
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("bad unmarshal result: want: %#v, got: %#v", tt.want, got)
			}
		})
	}
}
