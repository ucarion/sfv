package sfv

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
)

func Unmarshal(s string, v interface{}) error {
	scan := scanner{s: s, i: 0}
	scan.skipSP()

	switch v := v.(type) {
	case *Item:
		item, err := parseItem(&scan)
		if err != nil {
			return err
		}

		*v = item
	case *List:
		list, err := parseList(&scan)
		if err != nil {
			return err
		}

		*v = append(*v, list...)
	case *Dictionary:
		dict, err := parseDictionary(&scan)
		if err != nil {
			return err
		}

		for _, k := range dict.Keys {
			if _, ok := v.Map[k]; !ok {
				v.Keys = append(v.Keys, k)
			}

			if v.Map == nil {
				v.Map = map[string]Member{}
			}

			v.Map[k] = dict.Map[k]
		}

	// If the user did not provide one of the builtin SFV types, then we are
	// going to bind SFV data to the user-supplied type. But SFV's grammar is
	// such that you need to know in advance whether you're parsing an item,
	// list, or dictionary.
	//
	// So as a first step, we need to determine the data we're parsing. Then, we
	// can bind the parsed structure onto the user's supplied type.
	//
	// Since we're already checking v.(type), let's see if the user supplied a
	// "primitive" type. These correspond to an SFV item.
	case *bool, *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16,
		*uint32, *uint64, *float32, *float64, *string, *[]byte:
		item, err := parseItem(&scan)
		if err != nil {
			return err
		}

		if err := bindItemToSink(item, v); err != nil {
			return err
		}
	default:
		// The user may have supplied a struct (corresponds to an item), a slice
		// (corresponds to a list), or a map (corresponds to a dictionary). For
		// these cases, we need to use reflection.
		return fmt.Errorf("unsupported type: %T", v)
	}

	scan.skipSP()

	if !scan.isEOF() {
		return scan.parseError("illegal trailing characters")
	}

	return nil
}

func parseDictionary(s *scanner) (Dictionary, error) {
	var out Dictionary

	for {
		b, err := s.peek()
		if err != nil {
			break
		}

		key, err := parseKey(s)
		if err != nil {
			return Dictionary{}, err
		}

		var member Member

		b, err = s.peek()
		if err == nil && b == '=' {
			s.mustNext()
			member, err = parseListMember(s)
			if err != nil {
				return Dictionary{}, err
			}
		} else {
			params, err := parseParameters(s)
			if err != nil {
				return Dictionary{}, err
			}

			member = Member{
				IsItem: true,
				Item: Item{
					BareItem: BareItem{Type: BareItemTypeBoolean, Boolean: true},
					Params:   params,
				},
			}
		}

		if out.Map == nil {
			out.Map = map[string]Member{}
		}

		if _, ok := out.Map[key]; !ok {
			out.Keys = append(out.Keys, key)
		}

		out.Map[key] = member

		s.skipOWS()

		if s.isEOF() {
			return out, nil
		}

		b, err = s.next()
		if err != nil {
			return Dictionary{}, err
		}

		if b != ',' {
			return Dictionary{}, s.parseError("dictionary members must be delimited by ','")
		}

		s.skipOWS()

		if b, err := s.peek(); err != nil || b == ',' {
			return Dictionary{}, s.parseError("illegal trailing ','")
		}
	}

	return out, nil
}

func parseList(s *scanner) ([]Member, error) {
	var out []Member
	for {
		if s.isEOF() {
			break
		}

		member, err := parseListMember(s)
		if err != nil {
			return nil, err
		}

		out = append(out, member)

		s.skipOWS()

		if s.isEOF() {
			break
		}

		b, err := s.next()
		if err != nil {
			return nil, err
		}

		if b != ',' {
			return nil, s.parseError("list members must be delimited by ','")
		}

		s.skipOWS()

		if s.isEOF() {
			return nil, s.parseError("illegal trailing ',")
		}
	}

	return out, nil
}

func parseListMember(s *scanner) (Member, error) {
	b, err := s.peek()
	if err != nil {
		return Member{}, err
	}

	if b == '(' {
		innerList, err := parseInnerList(s)
		if err != nil {
			return Member{}, err
		}

		return Member{IsItem: false, InnerList: innerList}, nil
	}

	item, err := parseItem(s)
	if err != nil {
		return Member{}, err
	}

	return Member{IsItem: true, Item: item}, nil
}

func parseInnerList(s *scanner) (InnerList, error) {
	b, err := s.next()
	if err != nil {
		return InnerList{}, err
	}

	if b != '(' {
		return InnerList{}, s.parseError("inner list must start with '('")
	}

	items := []Item{}
	for {
		b, err := s.peek()
		if err != nil {
			break
		}

		s.skipSP()

		if b, err = s.peek(); err == nil && b == ')' {
			s.mustNext()
			params, err := parseParameters(s)
			if err != nil {
				return InnerList{}, err
			}

			return InnerList{Items: items, Params: params}, nil
		}

		item, err := parseItem(s)
		if err != nil {
			return InnerList{}, err
		}

		items = append(items, item)

		b, err = s.peek()
		if err != nil {
			return InnerList{}, err
		}

		if b != ' ' && b != ')' {
			return InnerList{}, s.parseError("inner lists items must be separated by ' '")
		}
	}

	return InnerList{}, s.parseError("unterminated inner list")
}

func parseItem(s *scanner) (Item, error) {
	bareItem, err := parseBareItem(s)
	if err != nil {
		return Item{}, err
	}

	params, err := parseParameters(s)
	if err != nil {
		return Item{}, err
	}

	return Item{BareItem: bareItem, Params: params}, nil
}

func parseBareItem(s *scanner) (BareItem, error) {
	b, err := s.peek()
	if err != nil {
		return BareItem{}, err
	}

	switch {
	case b == '-' || isDigit(b):
		return parseNumber(s)
	case b == '"':
		return parseString(s)
	case b == '*' || isAlpha(b):
		return parseToken(s)
	case b == ':':
		return parseByteSequence(s)
	case b == '?':
		return parseBoolean(s)
	default:
		return BareItem{}, fmt.Errorf("invalid start of bare item")
	}
}

func parseParameters(s *scanner) (Params, error) {
	out := Params{Map: map[string]BareItem{}, Keys: []string{}}

	for {
		b, err := s.peek()
		if err != nil {
			break // not an error; eof can end parameters
		}

		if b != ';' {
			break
		}

		s.mustNext()
		s.skipSP()

		key, err := parseKey(s)
		if err != nil {
			return Params{}, err
		}

		var value BareItem
		b, err = s.peek()
		if err != nil || b != '=' {
			// not an error; this just means that the param doesn't have a
			// value, so we use the default value instead
			value = BareItem{Type: BareItemTypeBoolean, Boolean: true}
		} else {
			s.mustNext()
			value, err = parseBareItem(s)
			if err != nil {
				return Params{}, err
			}
		}

		if _, ok := out.Map[key]; !ok {
			// this is a new key, append it to the ordering
			out.Keys = append(out.Keys, key)
		}

		out.Map[key] = value
	}

	return out, nil
}

func parseKey(s *scanner) (string, error) {
	b, err := s.peek()
	if err != nil {
		return "", err
	}

	if b != '*' && !isLCAlpha(b) {
		return "", s.parseError("bad start of key")
	}

	var buf []byte
	for {
		b, err := s.peek()
		if err != nil {
			return string(buf), nil // not an error; eof can terminate keys without values
		}

		if b != '_' && b != '-' && b != '.' && b != '*' && !isLCAlpha(b) && !isDigit(b) {
			return string(buf), nil
		}

		buf = append(buf, b)
		s.mustNext()
	}
}

func parseBoolean(s *scanner) (BareItem, error) {
	b, err := s.next()
	if err != nil {
		return BareItem{}, err
	}

	if b != '?' {
		return BareItem{}, s.parseError("boolean must start with '?'")
	}

	b, err = s.next()
	if err != nil {
		return BareItem{}, err
	}

	switch b {
	case '0':
		return BareItem{Type: BareItemTypeBoolean, Boolean: false}, nil
	case '1':
		return BareItem{Type: BareItemTypeBoolean, Boolean: true}, nil
	default:
		return BareItem{}, s.parseError("boolean value must be '0' or '1'")
	}
}

func parseNumber(s *scanner) (BareItem, error) {
	isInt := true      // are we parsing an integer, as opposed to a decimal?
	isPos := true      // what is the sign of the number?
	numBuf := []byte{} // a buffer of digits to parse

	b, err := s.peek()
	if err != nil {
		return BareItem{}, err
	}

	if b == '-' {
		isPos = false
		s.mustNext()
	}

	// detect an "empty" integer
	if _, err := s.peek(); err != nil {
		return BareItem{}, err
	}

	for {
		b, err := s.peek()
		if err != nil {
			break // eof is a valid way to end a number
		}

		ok := true
		switch {
		case isDigit(b):
			numBuf = append(numBuf, b)
			s.mustNext()
		case b == '.':
			if isInt {
				if len(numBuf) > 12 {
					return BareItem{}, s.parseError("too many digits in number")
				}

				numBuf = append(numBuf, b)
				isInt = false
				s.mustNext()
			} else {
				ok = false
			}
		default:
			ok = false
		}

		if !ok {
			break
		}

		if isInt && len(numBuf) > 15 {
			return BareItem{}, s.parseError("too many digits in number")
		}

		if !isInt && len(numBuf) > 16 {
			return BareItem{}, s.parseError("too many digits in number")
		}
	}

	if isInt {
		i, err := strconv.Atoi(string(numBuf))
		if err != nil {
			return BareItem{}, err
		}

		if isPos {
			return BareItem{Type: BareItemTypeInteger, Integer: int64(i)}, nil
		}

		return BareItem{Type: BareItemTypeInteger, Integer: int64(-i)}, nil
	}

	if numBuf[len(numBuf)-1] == '.' {
		return BareItem{}, s.parseError("number cannot end in '.'")
	}

	if len(numBuf)-bytes.Index(numBuf, []byte{'.'}) > 4 {
		return BareItem{}, s.parseError("too much precision in fractional part of decimal")
	}

	n, err := strconv.ParseFloat(string(numBuf), 64)
	if err != nil {
		panic(err) // should be unreachable
	}

	if isPos {
		return BareItem{Type: BareItemTypeDecimal, Decimal: n}, nil
	}

	return BareItem{Type: BareItemTypeDecimal, Decimal: -n}, nil
}

func parseString(s *scanner) (BareItem, error) {
	b, err := s.next()
	if err != nil {
		return BareItem{}, err
	}

	if b != '"' {
		return BareItem{}, s.parseError("string must start with '\"'")
	}

	var buf []byte
	for {
		b, err := s.next()
		if err != nil {
			return BareItem{}, err
		}

		switch {
		case b == '\\':
			b, err := s.next()
			if err != nil {
				return BareItem{}, err
			}

			if b != '\\' && b != '"' {
				return BareItem{}, s.parseError("only '\\' and '\"' may be escaped")
			}

			buf = append(buf, b)
		case b == '"':
			return BareItem{Type: BareItemTypeString, String: string(buf)}, nil
		case b != ' ' && !isVisible(b):
			return BareItem{}, s.parseError("strings must contain only spaces or visible ascii")
		default:
			buf = append(buf, b)
		}
	}
}

func parseToken(s *scanner) (BareItem, error) {
	b, err := s.peek()
	if err != nil {
		return BareItem{}, err // tokens cannot be empty
	}

	if b != '*' && !isAlpha(b) {
		return BareItem{}, s.parseError("invalid start of token")
	}

	var buf []byte
	for {
		b, err := s.peek()
		if err != nil {
			return BareItem{Type: BareItemTypeToken, Token: string(buf)}, nil // not an error; tokens can end at any time
		}

		if b != ':' && b != '/' && !isTChar(b) {
			return BareItem{Type: BareItemTypeToken, Token: string(buf)}, nil
		}

		buf = append(buf, b)
		s.mustNext()
	}
}

func parseByteSequence(s *scanner) (BareItem, error) {
	b, err := s.next()
	if err != nil {
		return BareItem{}, err
	}

	if b != ':' {
		return BareItem{}, s.parseError("byte sequence must start with ':'")
	}

	var buf []byte
	for {
		b, err := s.next()
		if err != nil {
			return BareItem{}, err
		}

		if b == ':' {
			break
		}

		buf = append(buf, b)
	}

	bytes, err := base64.StdEncoding.DecodeString(string(buf))
	if err != nil {
		return BareItem{}, err
	}

	return BareItem{Type: BareItemTypeBinary, Binary: bytes}, nil
}

type ParseError struct {
	Offset int
	msg    string
}

func (pe ParseError) Error() string {
	return pe.msg
}
