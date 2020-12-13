package sfv

import (
	"encoding/base64"
	"fmt"
	"math"
	"strings"
)

func Marshal(v interface{}) (string, error) {
	var w strings.Builder

	// todo: other types
	if v, ok := v.(Item); ok {
		if err := marshalItem(&w, v); err != nil {
			return "", err
		}
	}

	return w.String(), nil
}

func marshalItem(w *strings.Builder, v Item) error {
	if err := marshalBareItem(w, v.Value); err != nil {
		return err
	}

	if err := marshalParams(w, v.Params); err != nil {
		return err
	}

	return nil
}

func marshalBareItem(w *strings.Builder, v interface{}) error {
	switch v := v.(type) {
	case float64:
		return marshalNumber(w, v)
	case string:
		return marshalString(w, v)
	case []byte:
		return marshalByteSequence(w, v)
	case bool:
		return marshalBoolean(w, v)
	}

	return nil // todo: make this be an error instead
}

func marshalNumber(w *strings.Builder, v float64) error {
	// todo: check range

	if v == math.Round(v) {
		// an integer
		fmt.Fprintf(w, "%d", int64(v))
		return nil
	}

	return nil
}

func marshalString(w *strings.Builder, v string) error {
	// todo: check all chars ascii
	fmt.Fprint(w, "\"")
	for _, c := range v {
		if c == '\\' || c == '"' {
			fmt.Fprintf(w, "\\%s", string(c))
		} else {
			fmt.Fprintf(w, "%s", string(c))
		}
	}
	fmt.Fprint(w, "\"")
	return nil
}

func marshalByteSequence(w *strings.Builder, v []byte) error {
	fmt.Fprintf(w, ":%s:", base64.StdEncoding.EncodeToString(v))
	return nil
}

func marshalBoolean(w *strings.Builder, v bool) error {
	n := 0
	if v {
		n = 1
	}

	fmt.Fprintf(w, "?%d", n)
	return nil
}

func marshalParams(w *strings.Builder, v Params) error {
	return nil
}
