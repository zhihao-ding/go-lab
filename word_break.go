package main

import (
	"fmt"
)

var word = "dogandcat"
var dicts = []string{"do", "d", "gand", "g", "and", "dog", "cat", "c", "at"}
var results = [][]string{}

var cpyCounter = 0

func IsContainsInDicts(str string) bool {
	for _, v := range dicts {
		if v == str {
			return true
		}
	}
	return false
}

func breaker(matched []string, pos int) {
	if pos >= len(word) {
		results = append(results, matched)
		return
	}
	match := []byte{}
	for i := pos; i < len(word); i++ {
		match = append(match, word[i])
		aword := string(match)
		if IsContainsInDicts(aword) {
			cpyMatched := make([]string, len(matched))
			copy(cpyMatched, matched)
			cpyMatched = append(cpyMatched, aword)
			cpyCounter += 1
			breaker(cpyMatched, i + 1)
		}
	}
}

func main() {
	fmt.Println(word)
	fmt.Println(dicts)

	breaker([]string{}, 0)

	fmt.Println()
	fmt.Println(results)

	fmt.Println("cpy times: ", cpyCounter)
}

