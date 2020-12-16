package sfv

import (
	"fmt"
	"reflect"
)

func bindItem(i Item, v reflect.Value) error {
	// First, try to detect a primitive type. In this case, we'll just directly
	// bind the bare item to v, and ignore the parameters.
	switch v.Interface().(type) {
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16,
		uint32, uint64, float32, float64, string, []byte:
		return bindBareItem(i.BareItem, v.Addr().Interface())
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("cannot marshal item into %s", v.Type())
	}

	if err := bindParams(i.Params, v); err != nil {
		return fmt.Errorf("bind params: %w", err)
	}

	for j := 0; j < v.NumField(); j++ {
		if _, ok := v.Type().Field(j).Tag.Lookup("sfv"); !ok {
			if err := bindBareItem(i.BareItem, v.Field(j).Addr().Interface()); err != nil {
				return fmt.Errorf("bind bare item: %w", err)
			}

			break
		}
	}

	return nil
}

func bindList(l List, v reflect.Value) error {
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("cannot marshal list into %s", v.Type())
	}

	for _, m := range l {
		// In psuedo-code, what we want to do is:
		//
		// var T t
		// bindItem(m.Item, &t) (or bindInnerList, depending on m.IsItem)
		// *l = append(*l, t)
		//
		// Where l is of type []T.

		item := reflect.New(v.Type().Elem())

		if m.IsItem {
			if err := bindItem(m.Item, item.Elem()); err != nil {
				return err
			}
		} else {
			if err := bindInnerList(m.InnerList, item.Elem()); err != nil {
				return err
			}
		}

		v.Set(reflect.Append(v, item.Elem()))
	}

	return nil
}

func bindInnerList(l InnerList, v reflect.Value) error {
	// If v is a struct, then look for the first untagged field. If that field
	// is a slice, then we'll operate on that, and then we'll also look for
	// params.
	if v.Kind() == reflect.Struct {
		if err := bindParams(l.Params, v); err != nil {
			return fmt.Errorf("bind params: %w", err)
		}

		// Find the first untagged field in v, and try to bind the innerList
		// items to that field. We'll do that by simply reassigning v to the
		// relevant field.
		for j := 0; j < v.NumField(); j++ {
			if _, ok := v.Type().Field(j).Tag.Lookup("sfv"); !ok {
				v = v.Field(j)
				break
			}
		}
	}

	if v.Kind() != reflect.Slice {
		return fmt.Errorf("cannot marshal inner list into %s", v.Type())
	}

	for _, i := range l.Items {
		// This code is similar to that in bindList, but this time the elements
		// are always going to be items.
		item := reflect.New(v.Type().Elem())
		if err := bindItem(i, item.Elem()); err != nil {
			return err
		}

		v.Set(reflect.Append(v, item.Elem()))
	}

	return nil
}

func bindDictionary(d Dictionary, v reflect.Value) error {
	if v.Kind() != reflect.Map {
		return fmt.Errorf("cannot marshal dictionary into %s", v.Type())
	}

	// The zero value of a map is nil, which you can't assign to. So we'll do,
	// using reflection more or less this:
	//
	// if v == nil { *v = make(map[K]V) }
	//
	// Where map[K]V is whatever type v is.
	if v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}

	for k, m := range d.Map {
		item := reflect.New(v.Type().Elem())

		if m.IsItem {
			if err := bindItem(m.Item, item.Elem()); err != nil {
				return err
			}
		} else {
			if err := bindInnerList(m.InnerList, item.Elem()); err != nil {
				return err
			}
		}

		v.SetMapIndex(reflect.ValueOf(k), item.Elem())
	}

	return nil
}

func bindParams(p Params, v reflect.Value) error {
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("cannot marshal params into %s", v.Type())
	}

	for i := 0; i < v.Type().NumField(); i++ {
		if paramName, ok := v.Type().Field(i).Tag.Lookup("sfv"); ok {
			if paramValue, ok := p.Map[paramName]; ok {
				if err := bindBareItem(paramValue, v.Field(i).Addr().Interface()); err != nil {
					return fmt.Errorf("%s: %w", paramName, err)
				}
			}
		}
	}

	return nil
}

func bindBareItem(i BareItem, v interface{}) error {
	switch v := v.(type) {
	case *bool:
		if i.Type == BareItemTypeBoolean {
			*v = i.Boolean
			return nil
		}
	case *int:
		if i.Type == BareItemTypeInteger {
			*v = int(i.Integer)
			return nil
		}
	case *int8:
		if i.Type == BareItemTypeInteger {
			*v = int8(i.Integer)
			return nil
		}
	case *int16:
		if i.Type == BareItemTypeInteger {
			*v = int16(i.Integer)
			return nil
		}
	case *int32:
		if i.Type == BareItemTypeInteger {
			*v = int32(i.Integer)
			return nil
		}
	case *int64:
		if i.Type == BareItemTypeInteger {
			*v = int64(i.Integer)
			return nil
		}
	case *uint:
		if i.Type == BareItemTypeInteger {
			*v = uint(i.Integer)
			return nil
		}
	case *uint8:
		if i.Type == BareItemTypeInteger {
			*v = uint8(i.Integer)
			return nil
		}
	case *uint16:
		if i.Type == BareItemTypeInteger {
			*v = uint16(i.Integer)
			return nil
		}
	case *uint32:
		if i.Type == BareItemTypeInteger {
			*v = uint32(i.Integer)
			return nil
		}
	case *uint64:
		if i.Type == BareItemTypeInteger {
			*v = uint64(i.Integer)
			return nil
		}
	case *float32:
		if i.Type == BareItemTypeDecimal {
			*v = float32(i.Decimal)
			return nil
		}
	case *float64:
		if i.Type == BareItemTypeDecimal {
			*v = i.Decimal
			return nil
		}
	case *string:
		if i.Type == BareItemTypeString {
			*v = i.String
			return nil
		}

		if i.Type == BareItemTypeToken {
			*v = i.Token
			return nil
		}
	case *[]byte:
		if i.Type == BareItemTypeBinary {
			*v = i.Binary
			return nil
		}
	}

	return fmt.Errorf("cannot marshal to %T from %s", v, i.Type)
}
