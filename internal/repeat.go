package internal

import "strings"

func RepeatAndJoin(word, sep string, count int) string {
	if count < 1 {
		return ""
	}

	if count == 1 {
		return word
	}

	var buffer strings.Builder
	buffer.Grow(len(word) + (count-1)*(len(word)+len(sep)))

	buffer.WriteString(word)
	for index := 1; index < count; index++ {
		buffer.WriteString(sep)
		buffer.WriteString(word)
	}

	return buffer.String()
}
