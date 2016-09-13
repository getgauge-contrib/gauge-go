package stepImpl

import (
	"strconv"
	"testing"

	"github.com/manuviswam/gauge-go/gauge"
	m "github.com/manuviswam/gauge-go/models"
)

var vowels map[rune]bool

var _ = gauge.Step("Vowels in English language are <vowels>.", func(vowelString string) {
	vowels = make(map[rune]bool, 0)
	for _, ch := range vowelString {
		vowels[ch] = true
	}
})

var _ = gauge.Step("Almost all words have vowels <table>", func(tbl *m.Table) {
	var t *testing.T
	for _, row := range tbl.Rows {
		word := row.Cells[0]
		expectedCount, err := strconv.Atoi(row.Cells[1])
		if err != nil {
			t.Errorf("Failed to parse string %s to integer", row.Cells[1])
		}
		actualCount := countVowels(word)
		if actualCount != expectedCount {
			t.Errorf("got: %d, want: %d", actualCount, expectedCount)
		}
	}
})

var _ = gauge.Step("The word <word> has <expectedCount> vowels.", func(word string, expected string) {
	var t *testing.T
	actualCount := countVowels(word)
	expectedCount, err := strconv.Atoi(expected)
	if err != nil {
		t.Errorf("Failed to parse string %s to integer", expected)
	}
	if actualCount != expectedCount {
		t.Errorf("got: %d, want: %d", actualCount, expectedCount)
	}
})

func countVowels(word string) int {
	vowelCount := 0
	for _, ch := range word {
		if _, ok := vowels[ch]; ok {
			vowelCount++
		}
	}
	return vowelCount
}
