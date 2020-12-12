package sfv

func isAlpha(b byte) bool {
	return (b >= 0x41 && b <= 0x5A) || (b >= 0x61 && b <= 0x7A)
}

func isDigit(b byte) bool {
	return b >= 0x30 && b <= 0x39
}

func isVisible(b byte) bool {
	return b >= 0x21 && b <= 0x7E
}

func isTChar(b byte) bool {
	return isAlpha(b) ||
		isDigit(b) ||
		b == '!' ||
		b == '#' ||
		b == '$' ||
		b == '%' ||
		b == '&' ||
		b == '\'' ||
		b == '*' ||
		b == '+' ||
		b == '-' ||
		b == '.' ||
		b == '^' ||
		b == '_' ||
		b == '`' ||
		b == '|' ||
		b == '~'
}

func isLCAlpha(b byte) bool {
	return b >= 0x61 && b <= 0x7A
}
