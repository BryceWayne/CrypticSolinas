package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type State struct {
	Counter  int
	Curve    string
	Attempts map[string]bool
}

type HashInfo struct {
	Phrase string
	Hash   string
}

var targetHashesMap map[string]bool
var mu sync.Mutex

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

func generateHash(phrase string, ch chan<- HashInfo) {
	h := sha1.New()
	h.Write([]byte(phrase))
	hash := hex.EncodeToString(h.Sum(nil))
	ch <- HashInfo{Phrase: phrase, Hash: hash}
}

var targetHashes = [...]string{
	"3045AE6FC8422F64ED579528D38120EAE12196D5",
	"BD71344799D5C7FCDC45B59FA3B9AB8F6A948BC5",
	"C49D360886E704936A6678E1139D26B7819F7E90",
	"A335926AA319A27A1D00896A6773A4827ACDAC73",
	"D09E8800291CB85396CC6717393284AAA0DA64BA",
}

func main() {
	state, err := LoadState()
	if err != nil {
		state = &State{Counter: 0, Curve: "NIST P-192", Attempts: make(map[string]bool)}
	}

	words, err := loadDictionary()
	if err != nil {
		fmt.Println("Error loading dictionary:", err)
		return
	}

	targetHashesMap = make(map[string]bool)
	for _, hash := range targetHashes {
		targetHashesMap[hash] = true
	}

	var wg sync.WaitGroup
	ch := make(chan HashInfo)

	candidatePhrases := generateCandidatePhrases()

	many := 10_000_000
	bar := progressbar.Default(int64(many))
	for i := 0; i < many; i++ {
		randomWordLength := 1 + i%10
		randomPhrase := generateRandomPhrase(words, randomWordLength)
		candidatePhrases = append(candidatePhrases, randomPhrase)
		candidatePhrases = append(candidatePhrases, randomPhrase+".")
		bar.Add(1)
	}

	bar = progressbar.Default(int64(len(candidatePhrases)))

	for _, phrase := range candidatePhrases {
		mu.Lock()
		if _, exists := state.Attempts[phrase]; exists {
			mu.Unlock()
			continue
		}
		state.Attempts[phrase] = true
		mu.Unlock()

		wg.Add(1)
		go func(phrase string) {
			defer wg.Done()
			generateHash(phrase, ch)
			bar.Add(1)
		}(phrase)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for info := range ch {
		if targetHashesMap[info.Hash] {
			fmt.Printf("Match found! Hash: %s\n", info.Hash)
			if err := SaveSeed(info.Phrase, info.Hash); err != nil {
				fmt.Printf("Error saving seed: %s\n", err)
			}
		}
	}

	if err := SaveState(state); err != nil {
		fmt.Printf("Error saving state: %s\n", err)
	}
}
