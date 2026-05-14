package main

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func loadDictionary() ([]string, error) {
	filename := "dictionary.txt"
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words, scanner.Err()
}

func capitalizeWord(word string) string {
	if len(word) == 0 {
		return word
	}
	if word[0] >= 'a' && word[0] <= 'z' {
		b := make([]byte, len(word))
		copy(b, word)
		b[0] = b[0] - ('a' - 'A')
		return string(b)
	}
	return word
}

func generateRandomPhrase(words []string, phraseLength int) string {
	if phraseLength <= 0 {
		return ""
	}

	// Pre-allocate length estimation
	estimatedLen := phraseLength * 8
	buf := make([]byte, 0, estimatedLen)

	for i := 0; i < phraseLength; i++ {
		randomIndex := rand.Intn(len(words))
		word := words[randomIndex]

		if i == 0 {
			word = capitalizeWord(word)
		}

		buf = append(buf, word...)

		if i < phraseLength-1 {
			buf = append(buf, ' ')
		}
	}
	return string(buf)
}