package micro

func GetOffsetString(offset int, s, prefix, postfix string) string {

	if offset == 0 {
		return ""
	}

	var str string

	for i := 0; i < offset; i++ {
		str += s
	}

	return prefix + str + postfix
}
