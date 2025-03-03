package main

import (
	"encoding/json"
	"maps"
	"os"
	"slices"
	"strings"
)

func main() {
	jsonWords := make(map[string]int)
	raw, err := os.ReadFile("words_dictionary.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(raw, &jsonWords)
	if err != nil {
		panic(err)
	}

	words := slices.Collect(maps.Keys(jsonWords))
	slices.Sort(words)
	wordsOut := strings.Join(words, "\n")
	os.WriteFile("words.txt", []byte(wordsOut), 0777)
}
