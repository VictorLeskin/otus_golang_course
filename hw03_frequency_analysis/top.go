package hw03frequencyanalysis

import (
	//	"fmt"
	"sort"
	"strings"
)

type wordStat struct {
	word  string
	count int
}

var TopSize int = 10

func Top10(text string) []string {

	var ret []string
	words := strings.Fields(text)
	wordCount := make(map[string]int)

	for _, r := range words {
		wordCount[r]++
	}

	if len(wordCount) < TopSize {
		return ret
	}

	uniqueWordsStat := make([]wordStat, 0, len(wordCount))
	for word, cnt := range wordCount {
		uniqueWordsStat = append(uniqueWordsStat, wordStat{word: word, count: cnt})
	}

	sort.Slice(uniqueWordsStat, func(i, j int) bool {
		if uniqueWordsStat[i].count == uniqueWordsStat[j].count {
			return uniqueWordsStat[i].word < uniqueWordsStat[j].word
		}
		return uniqueWordsStat[i].count > uniqueWordsStat[j].count
	})

	for i := 0; i < TopSize; i++ {
		ret = append(ret, uniqueWordsStat[i].word)
	}
	// fmt.Println(uniqueWordsStat)

	// return nil
	return ret
}
