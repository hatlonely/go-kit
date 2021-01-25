package wrap

func ErrCode(err error) string {
	if err == nil {
		return "OK"
	}
	return "Unknown"
}
