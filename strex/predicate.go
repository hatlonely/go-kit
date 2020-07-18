package strex

func IsUpper(ch uint8) bool {
	return ch >= 'A' && ch <= 'Z'
}

func IsLower(ch uint8) bool {
	return ch >= 'a' && ch <= 'z'
}

func IsDigit(ch uint8) bool {
	return ch >= '0' && ch <= '9'
}

func IsXdigit(ch uint8) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func IsAlpha(ch uint8) bool {
	return IsLower(ch) || IsUpper(ch)
}

func IsAlnum(ch uint8) bool {
	return IsAlpha(ch) || IsDigit(ch)
}

func All(str string, op func(uint8) bool) bool {
	for i := range str {
		if !op(str[i]) {
			return false
		}
	}

	return true
}

func Any(str string, op func(uint8) bool) bool {
	for i := range str {
		if op(str[i]) {
			return true
		}
	}

	return false
}
