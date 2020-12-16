package sfv

import (
	"encoding/base64"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func Marshal(v interface{}) (string, error) {
	var w strings.Builder

	switch v := v.(type) {
	case Item:
		if err := marshalItem(&w, v); err != nil {
			return "", err
		}
	case List:
		if err := marshalList(&w, v); err != nil {
			return "", err
		}
	case Dictionary:
		if err := marshalDictionary(&w, v); err != nil {
			return "", err
		}
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16,
		uint32, uint64, float32, float64, string, []byte:
		item, err := unbindItem(reflect.ValueOf(v))
		if err != nil {
			return "", err
		}

		if err := marshalItem(&w, item); err != nil {
			return "", err
		}
	default:
		val := reflect.ValueOf(v)

		switch val.Kind() {
		case reflect.Struct:
			item, err := unbindItem(val)
			if err != nil {
				return "", err
			}

			if err := marshalItem(&w, item); err != nil {
				return "", err
			}
		case reflect.Slice:
			list, err := unbindList(val)
			if err != nil {
				return "", err
			}

			if err := marshalList(&w, list); err != nil {
				return "", err
			}
		case reflect.Map:
			dict, err := unbindDictionary(val)
			if err != nil {
				return "", err
			}

			if err := marshalDictionary(&w, dict); err != nil {
				return "", err
			}
		default:
			return "", fmt.Errorf("unsupported type: %T", v)
		}
	}

	return w.String(), nil
}

func marshalItem(w *strings.Builder, v Item) error {
	if err := marshalBareItem(w, v.BareItem); err != nil {
		return err
	}

	if err := marshalParams(w, v.Params); err != nil {
		return err
	}

	return nil
}

func marshalList(w *strings.Builder, v []Member) error {
	for i, m := range v {
		if m.IsItem {
			if err := marshalItem(w, m.Item); err != nil {
				return err
			}
		} else {
			if err := marshalInnerList(w, m.InnerList); err != nil {
				return err
			}
		}

		if i != len(v)-1 {
			fmt.Fprint(w, ", ")
		}
	}

	return nil
}

func marshalDictionary(w *strings.Builder, v Dictionary) error {
	for i, k := range v.Keys {
		if err := marshalKey(w, k); err != nil {
			return err
		}

		if v.Map[k].IsItem && v.Map[k].Item.BareItem.Type == BareItemTypeBoolean && v.Map[k].Item.BareItem.Boolean == true {
			if err := marshalParams(w, v.Map[k].Item.Params); err != nil {
				return err
			}
		} else {
			fmt.Fprint(w, "=")

			if v.Map[k].IsItem {
				if err := marshalItem(w, v.Map[k].Item); err != nil {
					return err
				}
			} else {
				if err := marshalInnerList(w, v.Map[k].InnerList); err != nil {
					return err
				}
			}
		}

		if i != len(v.Keys)-1 {
			fmt.Fprint(w, ", ")
		}
	}

	return nil
}

func marshalInnerList(w *strings.Builder, v InnerList) error {
	fmt.Fprint(w, "(")

	for i, m := range v.Items {
		if err := marshalItem(w, m); err != nil {
			return err
		}

		if i != len(v.Items)-1 {
			fmt.Fprintf(w, " ")
		}
	}

	fmt.Fprint(w, ")")

	if err := marshalParams(w, v.Params); err != nil {
		return err
	}

	return nil
}

func marshalBareItem(w *strings.Builder, v BareItem) error {
	switch v.Type {
	case BareItemTypeDecimal:
		return marshalDecimal(w, v.Decimal)
	case BareItemTypeInteger:
		return marshalInteger(w, v.Integer)
	case BareItemTypeString:
		return marshalString(w, v.String)
	case BareItemTypeToken:
		return marshalToken(w, v.Token)
	case BareItemTypeBinary:
		return marshalByteSequence(w, v.Binary)
	case BareItemTypeBoolean:
		return marshalBoolean(w, v.Boolean)
	default:
		return fmt.Errorf("unsupported bare item type: %v", v)
	}
}

func marshalDecimal(w *strings.Builder, v float64) error {
	if int64(v) < -999_999_999_999 || int64(v) > 999_999_999_999 {
		return fmt.Errorf("decimal out of range: %v", v)
	}

	// limit to three digits of precision past the decimal
	v = math.RoundToEven(v*1000) / 1000

	s := strconv.FormatFloat(v, 'f', -1, 64)
	if !strings.ContainsRune(s, '.') {
		s += ".0"
	}

	fmt.Fprint(w, s)
	return nil
}

func marshalInteger(w *strings.Builder, v int64) error {
	if v < -999_999_999_999_999 || v > 999_999_999_999_999 {
		return fmt.Errorf("integer out of range: %v", v)
	}

	fmt.Fprintf(w, "%d", v)
	return nil
}

func marshalString(w *strings.Builder, v string) error {
	fmt.Fprint(w, "\"")
	for _, c := range v {
		if !isVisible(byte(c)) && c != ' ' {
			return fmt.Errorf("invalid char in string: %c", c)
		}

		if c == '\\' || c == '"' {
			fmt.Fprintf(w, "\\%s", string(c))
		} else {
			fmt.Fprintf(w, "%s", string(c))
		}
	}
	fmt.Fprint(w, "\"")
	return nil
}

func marshalToken(w *strings.Builder, v string) error {
	for i, c := range v {
		if i == 0 && !isAlpha(byte(c)) && c != '*' {
			return fmt.Errorf("invalid first char in token: %v", c)
		}

		if i != 0 && !isTChar(byte(c)) && c != ':' && c != '/' {
			return fmt.Errorf("invalid char in token: %v", c)
		}
	}

	fmt.Fprintf(w, "%s", string(v))
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
	for _, k := range v.Keys {
		fmt.Fprintf(w, ";")
		if err := marshalKey(w, k); err != nil {
			return err
		}

		if !v.Map[k].isBoolTrue() {
			fmt.Fprint(w, "=")
			if err := marshalBareItem(w, v.Map[k]); err != nil {
				return err
			}
		}
	}

	return nil
}

func marshalKey(w *strings.Builder, v string) error {
	for i, c := range v {
		if i == 0 && !isLCAlpha(byte(c)) && c != '*' {
			return fmt.Errorf("invalid first char in key: %c", c)
		}

		if !isLCAlpha(byte(c)) && !isDigit(byte(c)) && c != '_' && c != '-' && c != '.' && c != '*' {
			return fmt.Errorf("invalid char in key: %c", c)
		}
	}

	fmt.Fprintf(w, "%s", v)
	return nil
}
