package sfv_test

import (
	"fmt"

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
	// public {6 0 0   [] true}
	// max-age {1 604800 0   [] false}
	// immutable {6 0 0   [] true}
}

// func ExampleUnmarshal_custom_bare_item() {
// 	var data string
// 	fmt.Println(sfv.Unmarshal("text/html; charset=UTF-8", &data))
// 	fmt.Println(data)
// 	// Output:
// 	// <nil>
// 	// text/html
// }

// func ExampleUnmarshal_custom_item() {
// 	type contentType struct {
// 		MediaType string
// 		Charset   string `sfv:"charset"`
// 		Boundary  string `sfv:"boundary"`
// 	}

// 	var data contentType
// 	fmt.Println(sfv.Unmarshal("text/html; charset=UTF-8", &data))
// 	fmt.Println(data)
// 	// Output:
// 	// <nil>
// }

// func TestUnmarshal_custom_bare_types(t *testing.T) {
// 	testCases := []struct {
// 		In      string
// 		Type    reflect.Kind
// 		Bool    bool
// 		String  string
// 		Int     int
// 		Int8    int8
// 		Int16   int16
// 		Int32   int32
// 		Int64   int64
// 		Uint    uint
// 		Uint8   uint8
// 		Uint16  uint16
// 		Uint32  uint32
// 		Uint64  uint64
// 		Float32 float32
// 		Float64 float64
// 		Bytes   []byte
// 	}{
// 		{In: "?0", Type: reflect.Bool, Bool: false},
// 		{In: "?1", Type: reflect.Bool, Bool: true},
// 		{In: "?1;foo=bar", Type: reflect.Bool, Bool: true},
// 		{In: "foo", Type: reflect.String, String: "foo"},
// 		{In: "\"foo\"", Type: reflect.String, String: "foo"},
// 		{In: "3", Type: reflect.Int, Int: 3},
// 		{In: "3", Type: reflect.Int8, Int8: 3},
// 		{In: "3", Type: reflect.Int16, Int16: 3},
// 		{In: "3", Type: reflect.Int32, Int32: 3},
// 		{In: "3", Type: reflect.Int64, Int64: 3},
// 		{In: "3", Type: reflect.Uint, Uint: 3},
// 		{In: "3", Type: reflect.Uint8, Uint8: 3},
// 		{In: "3", Type: reflect.Uint16, Uint16: 3},
// 		{In: "3", Type: reflect.Uint32, Uint32: 3},
// 		{In: "3", Type: reflect.Uint64, Uint64: 3},
// 		{In: "3.14", Type: reflect.Float32, Float32: 3.14},
// 		{In: "3.14", Type: reflect.Float64, Float64: 3.14},
// 		{In: ":aGVsbG8=:", Type: reflect.Array, Bytes: []byte{'h', 'e', 'l', 'l', 'o'}},
// 	}

// 	for _, tt := range testCases {
// 		t.Run(fmt.Sprintf("%s %s", tt.Type, tt.In), func(t *testing.T) {
// 			var want interface{}
// 			var got interface{}
// 			var err error

// 			switch tt.Type {
// 			case reflect.Bool:
// 				want = tt.Bool
// 				var v bool
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			case reflect.String:
// 				want = tt.String
// 				var v string
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			case reflect.Int:
// 				want = tt.Int
// 				var v int
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			case reflect.Int8:
// 				want = tt.Int8
// 				var v int8
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			case reflect.Int16:
// 				want = tt.Int16
// 				var v int16
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			case reflect.Int32:
// 				want = tt.Int32
// 				var v int32
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			case reflect.Int64:
// 				want = tt.Int64
// 				var v int64
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			case reflect.Uint:
// 				want = tt.Uint
// 				var v uint
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			case reflect.Uint8:
// 				want = tt.Uint8
// 				var v uint8
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			case reflect.Uint16:
// 				want = tt.Uint16
// 				var v uint16
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			case reflect.Uint32:
// 				want = tt.Uint32
// 				var v uint32
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			case reflect.Uint64:
// 				want = tt.Uint64
// 				var v uint64
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			case reflect.Float32:
// 				want = tt.Float32
// 				var v float32
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			case reflect.Float64:
// 				want = tt.Float64
// 				var v float64
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			case reflect.Array:
// 				want = tt.Bytes
// 				var v []byte
// 				err = sfv.Unmarshal(tt.In, &v)
// 				got = v
// 			default:
// 				panic("bad tt.Type")
// 			}

// 			if err != nil {
// 				t.Errorf("err: %v", err)
// 				return
// 			}

// 			if !reflect.DeepEqual(want, got) {
// 				t.Errorf("bad unmarshal result: want: %#v, got: %#v", want, got)
// 			}
// 		})
// 	}
// }

// // func TestUnmarshal_custom_type_string(t *testing.T) {
// // 	var data string
// // 	if err := sfv.Unmarshal("\"foo\"; bar=baz", &data); err != nil {
// // 		t.Errorf("err: %v", err)
// // 	}

// // 	if data != "foo" {
// // 		t.Errorf("data != \"foo\": %v", data)
// // 	}
// // }

// // func TestUnmarshal_custom_type_int(t *testing.T) {
// // 	var data int
// // 	if err := sfv.Unmarshal("3; bar=baz", &data); err != nil {
// // 		t.Errorf("err: %v", err)
// // 	}

// // 	if data != 3 {
// // 		t.Errorf("data != \"foo\": %v", data)
// // 	}
// // }
