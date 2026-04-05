package main

import (
	"sort"
	"strings"
)

func canonical(word string) string {
	runes := []rune(word)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}

func FindAnagramSets(words []string) map[string][]string {
	type group struct {
		firstWord string
		unique    map[string]struct{}
	}

	groups := make(map[string]*group, len(words))
	order := make([]string, 0, len(words))

	for _, word := range words {
		lower := strings.ToLower(word)
		signature := canonical(lower)

		g, ok := groups[signature]
		if !ok {
			g = &group{
				firstWord: lower,
				unique:    make(map[string]struct{}),
			}
			groups[signature] = g
			order = append(order, signature)
		}
		g.unique[lower] = struct{}{}
	}

	result := make(map[string][]string)
	for _, signature := range order {
		g := groups[signature]
		if len(g.unique) < 2 {
			continue
		}

		setWords := make([]string, 0, len(g.unique))
		for word := range g.unique {
			setWords = append(setWords, word)
		}
		sort.Strings(setWords)
		result[g.firstWord] = setWords
	}

	return result
}

func main() {}
