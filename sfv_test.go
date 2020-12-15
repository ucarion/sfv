package sfv_test

import (
	"encoding/base32"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ucarion/sfv"
)

// var skippedTestCases = map[string]struct{}{
// 	// Supporting this test case would require that we try multiple decoders
// 	// when parsing base64, which is undesirable. We currently fail on this test
// 	// case, and it's marked as may_fail in the test suite.
// 	"bad paddding": {}, // [sic]
// }

// func TestMarshalUnmarshal_StandardSuite(t *testing.T) {
// 	testFiles, err := filepath.Glob("structured-field-tests/*.json")
// 	if err != nil {
// 		t.Errorf("glob test files: %v", err)
// 	}

// 	for _, testFile := range testFiles {
// 		file, err := os.Open(testFile)
// 		if err != nil {
// 			t.Errorf("open test file: %v", err)
// 		}

// 		// We need to enable UseNumber here because the test suite distinguishes
// 		// between "1" and "1.0" in JSON numbers. For instance, the test case
// 		// named "single item parameterized list" contains a bare item
// 		// represented in JSON as "1.0", and expects that to be serialized as
// 		// "1.0", not "1".
// 		//
// 		// To deal with this, decodeBareItem will handle json.Number by looking
// 		// for a "." in the string representation of the number.
// 		decoder := json.NewDecoder(file)
// 		decoder.UseNumber()

// 		var testCases []testCase
// 		if err := decoder.Decode(&testCases); err != nil {
// 			t.Errorf("parse test file: %v", err)
// 		}

// 		for _, tt := range testCases {
// 			t.Run(tt.Name, func(t *testing.T) {
// 				if _, ok := skippedTestCases[tt.Name]; ok {
// 					t.SkipNow()
// 				}

// 				t.Run("unmarshal", func(t *testing.T) {
// 					expected, actual, err := tt.verifyUnmarshal()

// 					if tt.MustFail && err == nil {
// 						t.Errorf("must fail, but didn't")
// 					}

// 					if !tt.MustFail && err != nil {
// 						t.Errorf("must not fail, but did: %v", err)
// 					}

// 					if err == nil && !reflect.DeepEqual(expected, actual) {
// 						t.Errorf("bad unmarshal: want: %#v, got: %#v", expected, actual)
// 					}
// 				})

// 				t.Run("marshal", func(t *testing.T) {
// 					if tt.Expected == nil {
// 						return // nothing to marshal
// 					}

// 					expected, actual, err := tt.verifyMarshal()

// 					if tt.MustFail && err == nil {
// 						t.Errorf("must fail, but didn't")
// 					}

// 					if !tt.MustFail && err != nil {
// 						t.Errorf("must not fail, but did: %v", err)
// 					}

// 					if err == nil && expected != actual {
// 						t.Errorf("bad marshal: want: %#v, got: %#v", expected, actual)
// 					}
// 				})
// 			})
// 		}
// 	}
// }

// func TestMarshal_Serialization(t *testing.T) {
// 	testFiles, err := filepath.Glob("structured-field-tests/serialization-tests/*.json")
// 	if err != nil {
// 		t.Errorf("glob test files: %v", err)
// 	}

// 	for _, testFile := range testFiles {
// 		file, err := os.Open(testFile)
// 		if err != nil {
// 			t.Errorf("open test file: %v", err)
// 		}

// 		// We need to enable UseNumber here because the test suite distinguishes
// 		// between "1" and "1.0" in JSON numbers. For instance, the test case
// 		// named "single item parameterized list" contains a bare item
// 		// represented in JSON as "1.0", and expects that to be serialized as
// 		// "1.0", not "1".
// 		//
// 		// To deal with this, decodeBareItem will handle json.Number by looking
// 		// for a "." in the string representation of the number.
// 		decoder := json.NewDecoder(file)
// 		decoder.UseNumber()

// 		var testCases []testCase
// 		if err := decoder.Decode(&testCases); err != nil {
// 			t.Errorf("parse test file: %v", err)
// 		}

// 		for _, tt := range testCases {
// 			t.Run(tt.Name, func(t *testing.T) {
// 				if _, ok := skippedTestCases[tt.Name]; ok {
// 					t.SkipNow()
// 				}

// 				expected, actual, err := tt.verifyMarshal()

// 				if tt.MustFail && err == nil {
// 					t.Errorf("must fail, but didn't")
// 				}

// 				if !tt.MustFail && err != nil {
// 					t.Errorf("must not fail, but did: %v", err)
// 				}

// 				if err == nil && expected != actual {
// 					t.Errorf("bad marshal: want: %#v, got: %#v", expected, actual)
// 				}
// 			})
// 		}
// 	}
// }

// func (tt testCase) verifyUnmarshal() (interface{}, interface{}, error) {
// 	switch tt.HeaderType {
// 	case "item":
// 		var out sfv.Item
// 		for _, s := range tt.Raw {
// 			if err := sfv.Unmarshal(s, &out); err != nil {
// 				return nil, nil, err
// 			}
// 		}

// 		return decodeItem(tt.Expected), out, nil
// 	case "list":
// 		out := []sfv.Member{}
// 		for _, s := range tt.Raw {
// 			if err := sfv.Unmarshal(s, &out); err != nil {
// 				return nil, nil, err
// 			}
// 		}

// 		return decodeList(tt.Expected), out, nil
// 	case "dictionary":
// 		out := sfv.Dictionary{Map: map[string]sfv.Member{}, Keys: []string{}}
// 		for _, s := range tt.Raw {
// 			if err := sfv.Unmarshal(s, &out); err != nil {
// 				return nil, nil, err
// 			}
// 		}

// 		return decodeDictionary(tt.Expected), out, nil
// 	default:
// 		return nil, nil, fmt.Errorf("bad header type: %v", tt.HeaderType)
// 	}
// }

// func (tt testCase) verifyMarshal() (string, string, error) {
// 	expected := ""

// 	if len(tt.Raw) > 0 {
// 		expected = tt.Raw[0]
// 	}

// 	if len(tt.Canonical) > 0 {
// 		expected = tt.Canonical[0]
// 	}

// 	switch tt.HeaderType {
// 	case "item":
// 		in := decodeItem(tt.Expected)
// 		out, err := sfv.Marshal(in)
// 		if err != nil {
// 			return "", "", err
// 		}

// 		return expected, out, nil
// 	case "list":
// 		in := decodeList(tt.Expected)
// 		out, err := sfv.Marshal(in)
// 		if err != nil {
// 			return "", "", err
// 		}

// 		return expected, out, nil
// 	case "dictionary":
// 		in := decodeDictionary(tt.Expected)
// 		out, err := sfv.Marshal(in)
// 		if err != nil {
// 			return "", "", err
// 		}

// 		return expected, out, nil
// 	}

// 	return "", "", nil
// }

type testCase struct {
	Name       string      `json:"name"`
	Raw        []string    `json:"raw"`
	HeaderType string      `json:"header_type"`
	Expected   interface{} `json:"expected"`
	MustFail   bool        `json:"must_fail"`
	CanFail    bool        `json:"can_fail"`
	Canonical  []string    `json:"canonical"`
}

func TestUnmarshal_StdTestSuite(t *testing.T) {
	testCases, err := parseTestFiles("structured-field-tests/*.json")
	if err != nil {
		t.Errorf("parse test files: %v", err)
		return
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			// Supporting this test case would require that we try multiple
			// decoders when parsing base64, which is undesirable. We currently
			// fail on this test case, and it's marked as may_fail in the test
			// suite.
			if tt.Name == "bad paddding" { // [sic]
				return
			}

			actual, err := unmarshalTestCase(tt)

			if tt.MustFail {
				if err == nil {
					t.Errorf("must fail, but err is nil")
				}
			} else {
				if err != nil {
					t.Errorf("err: %v", err)
					return
				}

				expected := decodeTestCase(tt)
				if !reflect.DeepEqual(actual, expected) {
					t.Errorf("actual != expected: want: %#v, got: %#v", expected, actual)
					return
				}
			}
		})
	}
}

func TestMarshal_StdTestSuite(t *testing.T) {
	testCases, err := parseTestFiles("structured-field-tests/*.json")
	if err != nil {
		t.Errorf("parse test files: %v", err)
		return
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			// If there is no data to marshal, then there is nothing to test
			// here.
			if tt.Expected == nil {
				return
			}

			actual, err := marshalTestCase(tt)

			if tt.MustFail {
				if err == nil {
					t.Errorf("must fail, but err is nil")
				}
			} else {
				if err != nil {
					t.Errorf("err: %v", err)
					return
				}

				expected := tt.Raw[0]
				if len(tt.Canonical) > 0 {
					expected = tt.Canonical[0]
				}

				if !reflect.DeepEqual(actual, expected) {
					t.Errorf("actual != expected: want: %#v, got: %#v", expected, actual)
					return
				}
			}
		})
	}
}

func TestMarshal_StdSerializationTestSuite(t *testing.T) {
	testCases, err := parseTestFiles("structured-field-tests/serialisation-tests/*.json")
	if err != nil {
		t.Errorf("parse test files: %v", err)
		return
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			actual, err := marshalTestCase(tt)

			if tt.MustFail {
				if err == nil {
					t.Errorf("must fail, but err is nil")
				}
			} else {
				if err != nil {
					t.Errorf("err: %v", err)
					return
				}

				expected := ""

				if len(tt.Raw) > 0 {
					expected = tt.Raw[0]
				}

				if len(tt.Canonical) > 0 {
					expected = tt.Canonical[0]
				}

				if !reflect.DeepEqual(actual, expected) {
					t.Errorf("actual != expected: want: %#v, got: %#v", expected, actual)
					return
				}
			}
		})
	}
}

func unmarshalTestCase(tt testCase) (interface{}, error) {
	switch tt.HeaderType {
	case "item":
		var out sfv.Item
		for _, r := range tt.Raw {
			if err := sfv.Unmarshal(r, &out); err != nil {
				return nil, err
			}
		}

		return out, nil
	case "list":
		var out []sfv.Member
		for _, r := range tt.Raw {
			if err := sfv.Unmarshal(r, &out); err != nil {
				return nil, err
			}
		}

		return out, nil
	case "dictionary":
		var out sfv.Dictionary
		for _, r := range tt.Raw {
			if err := sfv.Unmarshal(r, &out); err != nil {
				return nil, err
			}
		}

		return out, nil
	default:
		panic("unknown header type")
	}
}

func marshalTestCase(tt testCase) (string, error) {
	switch tt.HeaderType {
	case "item":
		return sfv.Marshal(decodeItem(tt.Expected))
	case "list":
		return sfv.Marshal(decodeList(tt.Expected))
	case "dictionary":
		return sfv.Marshal(decodeDictionary(tt.Expected))
	default:
		panic("unknown header type")
	}
}

func decodeTestCase(tt testCase) interface{} {
	switch tt.HeaderType {
	case "item":
		return decodeItem(tt.Expected)
	case "list":
		return decodeList(tt.Expected)
	case "dictionary":
		return decodeDictionary(tt.Expected)
	default:
		panic("unknown error type")
	}
}

func parseTestFiles(glob string) ([]testCase, error) {
	paths, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}

	var out []testCase
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		// We need to enable UseNumber here because the test suite distinguishes
		// between "1" and "1.0" in JSON numbers. For instance, the test case
		// named "single item parameterized list" contains a bare item
		// represented in JSON as "1.0", and expects that to be serialized as
		// "1.0", not "1".
		//
		// To deal with this, decodeBareItem will handle json.Number by looking
		// for a "." in the string representation of the number.
		decoder := json.NewDecoder(file)
		decoder.UseNumber()

		var cases []testCase
		if err := decoder.Decode(&cases); err != nil {
			return nil, err
		}

		out = append(out, cases...)
	}

	return out, nil
}

func decodeDictionary(v interface{}) sfv.Dictionary {
	var out sfv.Dictionary
	for _, pair := range v.([]interface{}) {
		p := pair.([]interface{})
		k, v := p[0].(string), decodeMember(p[1])

		if out.Map == nil {
			out.Map = map[string]sfv.Member{}
		}

		out.Map[k] = v
		out.Keys = append(out.Keys, k)
	}

	return out
}

func decodeList(v interface{}) []sfv.Member {
	if v == nil {
		return nil
	}

	var out []sfv.Member
	for _, i := range v.([]interface{}) {
		out = append(out, decodeMember(i))
	}

	return out
}

func decodeMember(v interface{}) sfv.Member {
	// Inner lists are arrays of 2-elem arrays; items are 2-elem arrays where
	// the first element is never an array.
	//
	// So if i[0] is an array, we are dealing with an inner-list.
	arr := v.([]interface{})

	if _, ok := arr[0].([]interface{}); ok {
		return sfv.Member{IsItem: false, InnerList: decodeInnerList(v)}
	}

	return sfv.Member{IsItem: true, Item: decodeItem(v)}
}

func decodeInnerList(v interface{}) sfv.InnerList {
	arr := v.([]interface{})

	items := []sfv.Item{}
	for _, i := range arr[0].([]interface{}) {
		items = append(items, decodeItem(i))
	}

	return sfv.InnerList{Items: items, Params: decodeParams(arr[1])}
}

func decodeItem(v interface{}) sfv.Item {
	if v == nil {
		return sfv.Item{}
	}

	arr := v.([]interface{})
	return sfv.Item{BareItem: decodeBareItem(arr[0]), Params: decodeParams(arr[1])}
}

func decodeParams(v interface{}) sfv.Params {
	out := sfv.Params{Map: map[string]sfv.BareItem{}, Keys: []string{}}
	for _, pair := range v.([]interface{}) {
		p := pair.([]interface{})
		k, v := p[0].(string), decodeBareItem(p[1])

		out.Map[k] = v
		out.Keys = append(out.Keys, k)
	}

	return out
}

func decodeBareItem(v interface{}) sfv.BareItem {
	switch v := v.(type) {
	case bool:
		return sfv.BareItem{Type: sfv.BareItemTypeBoolean, Boolean: v}
	case string:
		return sfv.BareItem{Type: sfv.BareItemTypeString, String: v}
	case json.Number:
		if strings.ContainsRune(v.String(), '.') {
			n, err := v.Float64()
			if err != nil {
				panic(err)
			}

			return sfv.BareItem{Type: sfv.BareItemTypeDecimal, Decimal: n}
		}

		n, err := v.Int64()
		if err != nil {
			panic(err)
		}

		return sfv.BareItem{Type: sfv.BareItemTypeInteger, Integer: n}
	case map[string]interface{}:
		t := v["__type"]
		switch t {
		case "token":
			return sfv.BareItem{Type: sfv.BareItemTypeToken, Token: v["value"].(string)}
		case "binary":
			b, err := base32.StdEncoding.DecodeString(v["value"].(string))
			if err != nil {
				panic(err)
			}

			return sfv.BareItem{Type: sfv.BareItemTypeBinary, Binary: b}
		default:
			panic("bad __type")
		}
	default:
		panic("bad bare item")
	}
}
