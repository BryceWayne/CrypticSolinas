package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type State struct {
	Counter  int
	Curve    string
	Attempts map[string]bool // Using a map for fast lookup
}

var mu sync.Mutex // Mutex for concurrent map access

func SaveState(state *State) error {
	file, err := os.Create("state.json")
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(state, "", "    ") // Indentation of 4 spaces
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}

func SaveSeed(phrase, hash string) error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile("seed.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data := map[string]string{"Phrase": phrase, "Hash": hash}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	return err
}

func LoadState() (*State, error) {
	file, err := os.Open("state.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	state := &State{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(state)
	return state, err
}

func generateHash(phrase string, ch chan<- string, phraseHashes map[string]string) {
	h := sha1.New()
	h.Write([]byte(phrase))
	hash := hex.EncodeToString(h.Sum(nil))

	mu.Lock()
	phraseHashes[phrase] = hash
	mu.Unlock()

	ch <- hash
}

var targetHashes = [...]string{
	"3045AE6FC8422F64ED579528D38120EAE12196D5",
	"BD71344799D5C7FCDC45B59FA3B9AB8F6A948BC5",
	"C49D360886E704936A6678E1139D26B7819F7E90",
	"A335926AA319A27A1D00896A6773A4827ACDAC73",
	"D09E8800291CB85396CC6717393284AAA0DA64BA",
}

var curveNames = [...]string{
	"P-192", "P-224", "P-256", "P-384", "P-521",
}

func generateCandidatePhrases() []string {
	candidates := []string{}
	separators := []string{" ", ".", "(", ")", "[", "]", "{", "}", "Counter:", "Curve:", "Count:"}

	for i := 0; i < 2400; i++ {
		counterStr := strconv.Itoa(i)
		for _, curve := range curveNames {
			parts := []string{"Jerry", curve, counterStr}
			for _, sep1 := range separators {
				for _, sep2 := range separators {
					for _, sep3 := range separators {
						candidate := parts[0] + sep1 + parts[1] + sep2 + parts[2] + sep3
						candidates = append(candidates, candidate)
					}
				}
			}
		}
	}
	return candidates
}

func main() {
	// Load existing state
	state, err := LoadState()
	if err != nil {
		// Initialize if loading failed
		state = &State{Counter: 0, Curve: "NIST P-192", Attempts: make(map[string]bool)}
	}

	if state.Attempts == nil {
		state.Attempts = make(map[string]bool)
	}

	var wg sync.WaitGroup
	ch := make(chan string)

	candidatePhrases := generateCandidatePhrases()
	phraseHashes := make(map[string]string) // Create a map to store phrase and its corresponding hash

	bar := progressbar.Default(int64(len(candidatePhrases)))
	for _, phrase := range candidatePhrases {
		mu.Lock()
		if _, exists := state.Attempts[phrase]; exists {
			mu.Unlock()
			continue
		}
		state.Attempts[phrase] = true
		mu.Unlock()

		wg.Add(1)
		bar.Add(1)
		go generateHash(phrase, ch, phraseHashes) // Updated to populate phraseHashes
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	// Check hash in the main thread
	for hash := range ch {
		mu.Lock()
		for phrase, calculatedHash := range phraseHashes {
			if calculatedHash == hash {
				for _, target := range targetHashes {
					if hash == target {
						fmt.Printf("Match found! Hash: %s\n", hash)
						if err := SaveSeed(phrase, hash); err != nil {
							fmt.Printf("Error saving seed: %s\n", err)
						}
					}
				}
			}
		}
		mu.Unlock()
	}

	if err := SaveState(state); err != nil {
		fmt.Printf("Error saving state: %s\n", err)
	}
}
