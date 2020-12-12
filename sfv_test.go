package sfv_test

import (
	"encoding/base32"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ucarion/sfv"
)

type testCase struct {
	Name       string      `json:"name"`
	Raw        []string    `json:"raw"`
	HeaderType string      `json:"header_type"`
	Expected   interface{} `json:"expected"`
	MustFail   bool        `json:"must_fail"`
	CanFail    bool        `json:"can_fail"`
	Canonical  []string    `json:"canonical"`
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
				if tt.MustFail || tt.CanFail {
					t.SkipNow()
				}

				if tt.HeaderType != "item" {
					t.SkipNow()
				}

				switch tt.HeaderType {
				case "item":
					var out sfv.Item
					for _, s := range tt.Raw {
						if err := sfv.Unmarshal(s, &out); err != nil {
							t.Errorf("unmarshal input: %v", err)
						}
					}

					expected := decodeItem(tt.Expected)
					if !reflect.DeepEqual(out, expected) {
						t.Errorf("out != expected: want: %v, got: %v", expected, out)
					}
				}
			})
		}
	}
}

func decodeItem(v interface{}) sfv.Item {
	if v == nil {
		return sfv.Item{}
	}

	arr := v.([]interface{})

	params := sfv.Params{Map: map[string]interface{}{}, Keys: []string{}}
	for _, pair := range arr[1].([]interface{}) {
		p := pair.([]interface{})
		k, v := p[0].(string), decodeBareItem(p[1])

		params.Map[k] = v
		params.Keys = append(params.Keys, k)
	}

	return sfv.Item{Value: decodeBareItem(arr[0]), Params: params}
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
