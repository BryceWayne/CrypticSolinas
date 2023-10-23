package main

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"
)

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

func generateRandomPhrase(words []string, phraseLength int) string {
	rand.Seed(time.Now().UnixNano())
	phrase := ""
	for i := 0; i < phraseLength; i++ {
		randomIndex := rand.Intn(len(words))
		phrase += words[randomIndex]
		if i == 0 {
			phrase = strings.Title(phrase)
		}
		if i < phraseLength-1 {
			phrase += " "
		}
	}
	return phrase
}

// func main() {
// 	words, err := loadDictionary()
// 	if err != nil {
// 		fmt.Println("Error loading dictionary:", err)
// 		return
// 	}

// 	randomPhrase := generateRandomPhrase(words, 5) // Generate a 5-word random phrase
// 	fmt.Println("Generated random phrase:", randomPhrase)
// }
