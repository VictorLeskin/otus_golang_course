package hw03frequencyanalysis

import (
	//	"fmt"
	"sort"
	"strings"
	"unicode"
)

type wordStat struct {
	word  string
	count int
}

var TopSize = 10

func IsLine(s string) bool {
	for _, char := range s {
		if char != '-' {
			return false
		}
	}
	return true
}

func Trasform(s string) string {
	if IsLine(s) {
		if len(s) == 1 {
			return ""
		}
		return s
	}

	s = strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	})
	s = strings.ToLower(s)
	return s
}

func Top10(text string) []string {
	var ret []string
	words0 := strings.Fields(text)

	words := make([]string, 0, len(words0))
	for _, r := range words0 {
		if r1 := Trasform(r); r1 != "" {
			words = append(words, r1)
		}
	}

	wordCount := make(map[string]int)

	for _, r := range words {
		wordCount[r]++
	}

	if len(wordCount) == 0 {
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

	cnt := min(len(uniqueWordsStat), TopSize)
	for i := 0; i < cnt; i++ {
		ret = append(ret, uniqueWordsStat[i].word)
	}
	// fmt.Println(uniqueWordsStat)

	// return nil
	return ret
}
