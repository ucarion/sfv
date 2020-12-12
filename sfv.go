package sfv

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

type Item struct {
	Value  interface{}
	Params Params
}

type Params struct {
	Map  map[string]interface{}
	Keys []string
}

func Unmarshal(s string, v interface{}) error {
	scan := scanner{s: s, i: 0}

	scan.skipSP()

	// todo: support other types
	if v, ok := v.(*Item); ok {
		value, err := parseBareItem(&scan)
		if err != nil {
			return err
		}

		params, err := parseParameters(&scan)
		if err != nil {
			return err
		}

		*v = Item{Value: value, Params: params}
	}

	// return fmt.Errorf("bad type: %v", v)
	return nil
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

		s.next() // can't fail; we already peeked this
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
			s.next()
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
		s.next()
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
		s.next() // can't fail, already peeked
	}

	// detect an "empty" integer
	if _, err := s.peek(); err != nil {
		return 0, err
	}

	for {
		b, err = s.peek()
		if err != nil {
			break // eof is a valid way to end a number
		}

		// fmt.Println("loop", numBuf, string(b), s, err)

		ok := true
		switch {
		case isDigit(b):
			numBuf = append(numBuf, b)
			s.next() // can't fail, peeked either before or at end of loop
		case b == '.':
			if isInt {
				if len(numBuf) > 12 {
					return 0, s.parseError("too many digits in number")
				}

				numBuf = append(numBuf, b)
				isInt = false
				s.next() // can't fail, peeked either before or at end of loop
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

		b, err = s.peek()
	}

	if isInt {
		i, err := strconv.Atoi(string(numBuf))
		if err != nil {
			panic(err) // should be unreachable
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
			return "", s.parseError("unexpected end of string")
		}

		switch {
		case b == '\\':
			b, err := s.next()
			if err != nil {
				return "", s.parseError("unexpected end of string")
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
		s.next() // can't fail; we already peeked this
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

type scanner struct {
	s string
	i int
}

func (s *scanner) peek() (byte, error) {
	if s.i == len(s.s) {
		return 0, s.parseError("unexpected end of SFV input")
	}

	return s.s[s.i], nil
}

func (s *scanner) next() (byte, error) {
	b, err := s.peek()
	if err != nil {
		return 0, err
	}

	s.i++
	return b, nil
}

func (s scanner) parseError(msg string) error {
	return ParseError{Offset: s.i, msg: msg}
}

func (s *scanner) skipSP() {
	for {
		b, err := s.peek()
		if err != nil || b != ' ' {
			break
		}

		s.next()
	}
}

type ParseError struct {
	Offset int
	msg    string
}

func (pe ParseError) Error() string {
	return pe.msg
}
