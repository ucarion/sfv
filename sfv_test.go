package sfv_test

import (
	"encoding/base32"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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

		var testCases []testCase
		if err := json.NewDecoder(file).Decode(&testCases); err != nil {
			t.Errorf("parse test file: %v", err)
		}

		for _, tt := range testCases {
			t.Run(tt.Name, func(t *testing.T) {
				if _, ok := skippedTestCases[tt.Name]; ok {
					t.SkipNow()
				}

				t.Run("unmarshal", func(t *testing.T) {
					err := tt.verifyUnmarshal()

					if tt.MustFail && err == nil {
						t.Errorf("must fail, but didn't")
					}

					if !tt.MustFail && err != nil {
						t.Errorf("must not fail, but did: %v", err)
					}
				})

				t.Run("marshal", func(t *testing.T) {
					err := tt.verifyMarshal()

					if tt.MustFail && err == nil {
						t.Errorf("must fail, but didn't")
					}

					if !tt.MustFail && err != nil {
						t.Errorf("must not fail, but did: %v", err)
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

func (tt testCase) verifyUnmarshal() error {
	switch tt.HeaderType {
	case "item":
		var out sfv.Item
		for _, s := range tt.Raw {
			if err := sfv.Unmarshal(s, &out); err != nil {
				return fmt.Errorf("unmarshal input: %w", err)
			}
		}

		expected := decodeItem(tt.Expected)
		if !reflect.DeepEqual(out, expected) {
			return fmt.Errorf("out != expected: want: %#v, got: %#v", expected, out)
		}
	case "list":
		out := []sfv.Member{}
		for _, s := range tt.Raw {
			if err := sfv.Unmarshal(s, &out); err != nil {
				return fmt.Errorf("unmarshal input: %w", err)
			}
		}

		expected := decodeList(tt.Expected)
		if !reflect.DeepEqual(out, expected) {
			return fmt.Errorf("out != expected: want: %#v, got: %#v", expected, out)
		}
	case "dictionary":
		out := sfv.Dictionary{Map: map[string]sfv.Member{}, Keys: []string{}}
		for _, s := range tt.Raw {
			if err := sfv.Unmarshal(s, &out); err != nil {
				return fmt.Errorf("unmarshal input: %w", err)
			}
		}

		expected := decodeDictionary(tt.Expected)
		if !reflect.DeepEqual(out, expected) {
			return fmt.Errorf("out != expected: want: %#v, got: %#v", expected, out)
		}
	default:
		return fmt.Errorf("bad header type: %v", tt.HeaderType)
	}

	return nil
}

func (tt testCase) verifyMarshal() error {
	switch tt.HeaderType {
	case "item":
		in := decodeItem(tt.Expected)
		out, err := sfv.Marshal(in)
		if err != nil {
			return err
		}

		expected := tt.Raw[0]
		if len(tt.Canonical) > 0 {
			expected = tt.Canonical[0]
		}

		if out != expected {
			return fmt.Errorf("out != expected, want: %#v, got: %#v", expected, out)
		}
	}

	return nil
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
	case bool, float64, string:
		return v
	case map[string]interface{}:
		t := v["__type"]
		switch t {
		case "token":
			return v["value"]
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
