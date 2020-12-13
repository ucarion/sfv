package sfv

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

func Unmarshal(s string, v interface{}) error {
	scan := scanner{s: s, i: 0}

	scan.skipSP()

	// todo: support other types
	if v, ok := v.(*Item); ok {
		item, err := parseItem(&scan)
		if err != nil {
			return err
		}

		*v = item
	}

	if v, ok := v.(*[]Member); ok {
		list, err := parseList(&scan)
		if err != nil {
			return err
		}

		*v = append(*v, list...)
	}

	if v, ok := v.(*Dictionary); ok {
		dict, err := parseDictionary(&scan)
		if err != nil {
			return err
		}

		if v == nil {
			*v = dict
		}

		for _, k := range dict.Keys {
			if _, ok := v.Map[k]; !ok {
				v.Keys = append(v.Keys, k)
			}

			v.Map[k] = dict.Map[k]
		}
	}

	return nil
}

func parseDictionary(s *scanner) (Dictionary, error) {
	out := Dictionary{Map: map[string]Member{}, Keys: []string{}}

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

			member = Member{IsItem: true, Item: Item{Value: true, Params: params}}
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
	out := []Member{}

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
	value, err := parseBareItem(s)
	if err != nil {
		return Item{}, err
	}

	params, err := parseParameters(s)
	if err != nil {
		return Item{}, err
	}

	return Item{Value: value, Params: params}, nil
}

func parseBareItem(s *scanner) (interface{}, error) {
	b, err := s.peek()
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("invalid start of raw item")
	}
}

func parseParameters(s *scanner) (Params, error) {
	out := Params{Map: map[string]interface{}{}, Keys: []string{}}

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

		var value interface{}
		b, err = s.peek()
		if err != nil || b != '=' {
			// not an error; this just means that the param doesn't have a
			// value, so we use the default value instead
			value = true
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

func parseBoolean(s *scanner) (bool, error) {
	b, err := s.next()
	if err != nil {
		return false, err
	}

	if b != '?' {
		return false, s.parseError("boolean must start with '?'")
	}

	b, err = s.next()
	if err != nil {
		return false, err
	}

	switch b {
	case '0':
		return false, nil
	case '1':
		return true, nil
	default:
		return false, s.parseError("boolean value must be '0' or '1'")
	}
}

func parseNumber(s *scanner) (float64, error) {
	isInt := true      // are we parsing an integer, as opposed to a decimal?
	isPos := true      // what is the sign of the number?
	numBuf := []byte{} // a buffer of digits to parse

	b, err := s.peek()
	if err != nil {
		return 0, err
	}

	if b == '-' {
		isPos = false
		s.mustNext()
	}

	// detect an "empty" integer
	if _, err := s.peek(); err != nil {
		return 0, err
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
					return 0, s.parseError("too many digits in number")
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
			return 0, s.parseError("too many digits in number")
		}

		if !isInt && len(numBuf) > 16 {
			return 0, s.parseError("too many digits in number")
		}
	}

	if isInt {
		i, err := strconv.Atoi(string(numBuf))
		if err != nil {
			return 0, err
		}

		if isPos {
			return float64(i), nil
		}

		return float64(-i), nil
	}

	if numBuf[len(numBuf)-1] == '.' {
		return 0, s.parseError("number cannot end in '.'")
	}

	// todo: count chars past .

	n, err := strconv.ParseFloat(string(numBuf), 64)
	if err != nil {
		panic(err) // should be unreachable
	}

	if isPos {
		return n, nil
	}

	return -n, nil
}

func parseString(s *scanner) (string, error) {
	b, err := s.next()
	if err != nil {
		return "", err
	}

	if b != '"' {
		return "", s.parseError("string must start with '\"'")
	}

	var buf []byte
	for {
		b, err := s.next()
		if err != nil {
			return "", err
		}

		switch {
		case b == '\\':
			b, err := s.next()
			if err != nil {
				return "", err
			}

			if b != '\\' && b != '"' {
				return "", s.parseError("only '\\' and '\"' may be escaped")
			}

			buf = append(buf, b)
		case b == '"':
			return string(buf), nil
		case b != ' ' && !isVisible(b):
			return "", s.parseError("strings must contain only spaces or visible ascii")
		default:
			buf = append(buf, b)
		}
	}
}

func parseToken(s *scanner) (string, error) {
	b, err := s.peek()
	if err != nil {
		return "", err // tokens cannot be empty
	}

	if b != '*' && !isAlpha(b) {
		return "", s.parseError("invalid start of token")
	}

	var buf []byte
	for {
		b, err := s.peek()
		if err != nil {
			return string(buf), nil // not an error; tokens can end at any time
		}

		if b != ':' && b != '/' && !isTChar(b) {
			return string(buf), nil
		}

		buf = append(buf, b)
		s.mustNext()
	}
}

func parseByteSequence(s *scanner) ([]byte, error) {
	b, err := s.next()
	if err != nil {
		return nil, err
	}

	if b != ':' {
		return nil, s.parseError("byte sequence must start with ':'")
	}

	var buf []byte
	for {
		b, err := s.next()
		if err != nil {
			return nil, err
		}

		if b == ':' {
			break
		}

		buf = append(buf, b)
	}

	return base64.StdEncoding.DecodeString(string(buf))
}

type ParseError struct {
	Offset int
	msg    string
}

func (pe ParseError) Error() string {
	return pe.msg
}
