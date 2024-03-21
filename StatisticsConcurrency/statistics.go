package statistics

import "strings"

func WordCount(line string) int {
	wordSlice := strings.Fields(line)
	return len(wordSlice)
}
