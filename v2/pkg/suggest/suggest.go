package suggest

import (
	"sort"

	"github.com/agnivade/levenshtein"
)

type priorityWord struct {
	word     string
	priority int
}

func Suggest(words []string, target string, size int) []string {
	priorityWords := make([]priorityWord, 0, len(words))

	for _, word := range words {
		var distance int

		if len(word) > len(target) {
			distance = levenshtein.ComputeDistance(target, word[:len(target)])*10 +
				levenshtein.ComputeDistance(target, word[len(target):])
		} else {
			distance = levenshtein.ComputeDistance(target, word) * 10
		}

		priorityWords = append(priorityWords, priorityWord{word: word, priority: distance})
	}

	sort.Slice(priorityWords, func(i, j int) bool {
		return priorityWords[i].priority < priorityWords[j].priority
	})

	resultSize := len(priorityWords)
	if resultSize > size {
		resultSize = size
	}

	results := make([]string, 0, resultSize)
	for _, suggest := range priorityWords[:resultSize] {
		results = append(results, suggest.word)
	}

	return results
}
