package sfv

import (
	"fmt"
	"reflect"
)

func unbindItem(v reflect.Value) (Item, error) {
	switch v := v.Interface().(type) {
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16,
		uint32, uint64, float32, float64, string, []byte:
		bareItem, err := unbindBareItem(v)
		if err != nil {
			return Item{}, err
		}

		return Item{BareItem: bareItem}, nil
	}

	if v.Kind() != reflect.Struct {
		return Item{}, fmt.Errorf("cannot unmarshal to item from %s", v.Type())
	}

	var bareItem BareItem
	params := Params{Map: map[string]BareItem{}}

	for i := 0; i < v.NumField(); i++ {
		if tag, ok := v.Type().Field(i).Tag.Lookup("sfv"); ok {
			if v.Field(i).IsZero() {
				continue
			}

			bareItem, err := unbindBareItem(v.Field(i).Interface())
			if err != nil {
				return Item{}, err
			}

			params.Keys = append(params.Keys, tag)
			params.Map[tag] = bareItem
		} else {
			var err error
			bareItem, err = unbindBareItem(v.Field(i).Interface())
			if err != nil {
				return Item{}, err
			}
		}
	}

	return Item{BareItem: bareItem, Params: params}, nil
}

func unbindList(v reflect.Value) (List, error) {
	if v.Kind() != reflect.Slice {
		return List{}, fmt.Errorf("cannot unmarshal to list from %s", v.Type())
	}

	var out List
	for i := 0; i < v.Len(); i++ {
		member, err := unbindMember(v.Index(i))
		if err != nil {
			return List{}, err
		}

		out = append(out, member)
	}

	return out, nil
}

func unbindInnerList(v reflect.Value) (InnerList, error) {
	var params Params
	if v.Kind() == reflect.Struct {
		var err error
		params, err = unbindParams(v)
		if err != nil {
			return InnerList{}, err
		}

		for j := 0; j < v.NumField(); j++ {
			if _, ok := v.Type().Field(j).Tag.Lookup("sfv"); !ok {
				v = v.Field(j)
				break
			}
		}
	}

	if v.Kind() != reflect.Slice {
		return InnerList{}, fmt.Errorf("cannot unmarshal to inner list from %s", v.Type())
	}

	var items []Item
	for i := 0; i < v.Len(); i++ {
		item, err := unbindItem(v.Index(i))
		if err != nil {
			return InnerList{}, err
		}

		items = append(items, item)
	}

	return InnerList{Items: items, Params: params}, nil
}

func unbindDictionary(v reflect.Value) (Dictionary, error) {
	if v.Kind() != reflect.Map || v.Type().Key().Kind() != reflect.String {
		return Dictionary{}, fmt.Errorf("cannot unmarshal to map from %s", v.Type())
	}

	out := Dictionary{Map: map[string]Member{}}

	iter := v.MapRange()
	for iter.Next() {
		member, err := unbindMember(iter.Value())
		if err != nil {
			return Dictionary{}, err
		}

		key := iter.Key().Interface().(string)
		out.Keys = append(out.Keys, key)
		out.Map[key] = member
	}

	return out, nil
}

func unbindMember(v reflect.Value) (Member, error) {
	isInnerList := v.Type().Kind() == reflect.Slice
	if !isInnerList && v.Type().Kind() == reflect.Struct {
		for j := 0; j < v.Type().NumField(); j++ {
			if _, ok := v.Type().Field(j).Tag.Lookup("sfv"); !ok {
				isInnerList = v.Type().Field(j).Type.Kind() == reflect.Slice
				break
			}
		}
	}

	if isInnerList {
		// v is some sort of slice-of-slices, so we're doing inner lists.
		innerList, err := unbindInnerList(v)
		if err != nil {
			return Member{}, err
		}

		return Member{
			IsItem:    false,
			InnerList: innerList,
		}, nil
	}

	// We're doing items.
	item, err := unbindItem(v)
	if err != nil {
		return Member{}, err
	}

	return Member{
		IsItem: true,
		Item:   item,
	}, nil

}

func unbindParams(v reflect.Value) (Params, error) {
	if v.Kind() != reflect.Struct {
		return Params{}, fmt.Errorf("cannot unmarshal to params from %s", v.Type())
	}

	params := Params{Map: map[string]BareItem{}}

	for i := 0; i < v.NumField(); i++ {
		if tag, ok := v.Type().Field(i).Tag.Lookup("sfv"); ok {
			if v.Field(i).IsZero() {
				continue
			}

			bareItem, err := unbindBareItem(v.Field(i).Interface())
			if err != nil {
				return Params{}, err
			}

			params.Keys = append(params.Keys, tag)
			params.Map[tag] = bareItem
		}
	}

	return params, nil
}

func unbindBareItem(v interface{}) (BareItem, error) {
	switch v := v.(type) {
	case int:
		return BareItem{Type: BareItemTypeInteger, Integer: int64(v)}, nil
	case string:
		// todo: how to decide whether to use token or string?
		return BareItem{Type: BareItemTypeToken, Token: v}, nil
	case float64:
		return BareItem{Type: BareItemTypeDecimal, Decimal: v}, nil
	}

	return BareItem{}, fmt.Errorf("cannot unmarshal from %T", v)
}
