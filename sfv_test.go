package sfv_test

import (
	"encoding/base32"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ucarion/sfv"
)

var skippedTestCases = map[string]struct{}{
	// Supporting this test case would require that we try multiple decoders
	// when parsing base64, which is undesirable. We currently fail on this test
	// case, and it's marked as may_fail in the test suite.
	"bad paddding": {}, // [sic]
}

func TestParse_StandardSuite(t *testing.T) {
	testFiles, err := filepath.Glob("structured-field-tests/*.json")
	if err != nil {
		t.Errorf("glob test files: %v", err)
	}

	for _, testFile := range testFiles {
		file, err := os.Open(testFile)
		if err != nil {
			t.Errorf("open test file: %v", err)
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

		var testCases []testCase
		if err := decoder.Decode(&testCases); err != nil {
			t.Errorf("parse test file: %v", err)
		}

		for _, tt := range testCases {
			t.Run(tt.Name, func(t *testing.T) {
				if _, ok := skippedTestCases[tt.Name]; ok {
					t.SkipNow()
				}

				t.Run("unmarshal", func(t *testing.T) {
					expected, actual, err := tt.verifyUnmarshal()

					if tt.MustFail && err == nil {
						t.Errorf("must fail, but didn't")
					}

					if !tt.MustFail && err != nil {
						t.Errorf("must not fail, but did: %v", err)
					}

					if err == nil && !reflect.DeepEqual(expected, actual) {
						t.Errorf("bad unmarshal: want: %#v, got: %#v", expected, actual)
					}
				})

				t.Run("marshal", func(t *testing.T) {
					if tt.Expected == nil {
						return // nothing to marshal
					}

					expected, actual, err := tt.verifyMarshal()

					if tt.MustFail && err == nil {
						t.Errorf("must fail, but didn't")
					}

					if !tt.MustFail && err != nil {
						t.Errorf("must not fail, but did: %v", err)
					}

					if err == nil && expected != actual {
						t.Errorf("bad marshal: want: %#v, got: %#v", expected, actual)
					}
				})
			})
		}
	}
}

type testCase struct {
	Name       string      `json:"name"`
	Raw        []string    `json:"raw"`
	HeaderType string      `json:"header_type"`
	Expected   interface{} `json:"expected"`
	MustFail   bool        `json:"must_fail"`
	CanFail    bool        `json:"can_fail"`
	Canonical  []string    `json:"canonical"`
}

func (tt testCase) verifyUnmarshal() (interface{}, interface{}, error) {
	switch tt.HeaderType {
	case "item":
		var out sfv.Item
		for _, s := range tt.Raw {
			if err := sfv.Unmarshal(s, &out); err != nil {
				return nil, nil, err
			}
		}

		return decodeItem(tt.Expected), out, nil
	case "list":
		out := []sfv.Member{}
		for _, s := range tt.Raw {
			if err := sfv.Unmarshal(s, &out); err != nil {
				return nil, nil, err
			}
		}

		return decodeList(tt.Expected), out, nil
	case "dictionary":
		out := sfv.Dictionary{Map: map[string]sfv.Member{}, Keys: []string{}}
		for _, s := range tt.Raw {
			if err := sfv.Unmarshal(s, &out); err != nil {
				return nil, nil, err
			}
		}

		return decodeDictionary(tt.Expected), out, nil
	default:
		return nil, nil, fmt.Errorf("bad header type: %v", tt.HeaderType)
	}
}

func (tt testCase) verifyMarshal() (string, string, error) {
	expected := tt.Raw[0]
	if len(tt.Canonical) > 0 {
		expected = tt.Canonical[0]
	}

	switch tt.HeaderType {
	case "item":
		in := decodeItem(tt.Expected)
		out, err := sfv.Marshal(in)
		if err != nil {
			return "", "", err
		}

		return expected, out, nil
	case "list":
		in := decodeList(tt.Expected)
		out, err := sfv.Marshal(in)
		if err != nil {
			return "", "", err
		}

		return expected, out, nil
	}

	return "", "", nil
}

func decodeDictionary(v interface{}) sfv.Dictionary {
	out := sfv.Dictionary{Map: map[string]sfv.Member{}, Keys: []string{}}
	for _, pair := range v.([]interface{}) {
		p := pair.([]interface{})
		k, v := p[0].(string), decodeMember(p[1])

		out.Map[k] = v
		out.Keys = append(out.Keys, k)
	}

	return out
}

func decodeList(v interface{}) []sfv.Member {
	if v == nil {
		return nil
	}

	out := []sfv.Member{}
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
	return sfv.Item{Value: decodeBareItem(arr[0]), Params: decodeParams(arr[1])}
}

func decodeParams(v interface{}) sfv.Params {
	out := sfv.Params{Map: map[string]interface{}{}, Keys: []string{}}
	for _, pair := range v.([]interface{}) {
		p := pair.([]interface{})
		k, v := p[0].(string), decodeBareItem(p[1])

		out.Map[k] = v
		out.Keys = append(out.Keys, k)
	}

	return out
}

func decodeBareItem(v interface{}) interface{} {
	switch v := v.(type) {
	case bool, string:
		return v
	case json.Number:
		if strings.ContainsRune(v.String(), '.') {
			n, err := v.Float64()
			if err != nil {
				panic(err)
			}

			return n
		}

		n, err := v.Int64()
		if err != nil {
			panic(err)
		}

		return n
	case map[string]interface{}:
		t := v["__type"]
		switch t {
		case "token":
			return sfv.Token(v["value"].(string))
		case "binary":
			b, err := base32.StdEncoding.DecodeString(v["value"].(string))
			if err != nil {
				panic(err)
			}
			return b
		default:
			panic("bad __type")
		}
	default:
		panic("bad bare item")
	}
}
