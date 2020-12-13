package sfv

type scanner struct {
	s string
	i int
}

func (s *scanner) isEOF() bool {
	return s.i == len(s.s)
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

func (s *scanner) mustNext() {
	if _, err := s.next(); err != nil {
		panic(err)
	}
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

func (s *scanner) skipOWS() {
	for {
		b, err := s.peek()
		if err != nil || !isOWS(b) {
			break
		}

		s.next()
	}
}

func (s scanner) parseError(msg string) error {
	return ParseError{Offset: s.i, msg: msg}
}
