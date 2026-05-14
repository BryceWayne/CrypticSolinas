package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type HashInfo struct {
	Phrase string
	Hash   string
}

var targetHashesMap map[[20]byte]bool
var mu sync.Mutex

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

var targetHashes = [...]string{
	"3045AE6FC8422F64ED579528D38120EAE12196D5",
	"BD71344799D5C7FCDC45B59FA3B9AB8F6A948BC5",
	"C49D360886E704936A6678E1139D26B7819F7E90",
	"A335926AA319A27A1D00896A6773A4827ACDAC73",
	"D09E8800291CB85396CC6717393284AAA0DA64BA",
}

func main() {
	mode := flag.String("mode", "cpu", "Execution mode: cpu or gpu")
	threads := flag.Int("threads", runtime.NumCPU(), "Number of threads for CPU mode")
	flag.Parse()

	words, err := loadDictionary()
	if err != nil {
		fmt.Println("Error loading dictionary:", err)
		return
	}

	targetHashesMap = make(map[[20]byte]bool)
	for _, hashHex := range targetHashes {
		decoded, err := hex.DecodeString(hashHex)
		if err == nil && len(decoded) == 20 {
			var b [20]byte
			copy(b[:], decoded)
			targetHashesMap[b] = true
		}
	}

	candidateCh := make(chan string, 100000)

	// Start phrase generation
	var wgGen sync.WaitGroup
	wgGen.Add(1)
	go func() {
		defer wgGen.Done()
		generateCandidatePhrases(candidateCh)

		many := 10_000_000
		for i := 0; i < many; i++ {
			randomWordLength := 1 + i%10
			randomPhrase := generateRandomPhrase(words, randomWordLength)
			candidateCh <- randomPhrase
			candidateCh <- randomPhrase + "."
		}
		close(candidateCh)
	}()

	if *mode == "gpu" {
		fmt.Println("GPU mode selected. Dumping phrases to candidates.txt...")
		file, err := os.Create("candidates.txt")
		if err != nil {
			fmt.Println("Error creating candidates.txt:", err)
			return
		}
		defer file.Close()

		bar := progressbar.Default(-1, "Generating")
		for phrase := range candidateCh {
			file.WriteString(phrase + "\n")
			bar.Add(1)
		}
		wgGen.Wait()

		// Write target hashes to file for hashcat
		targetFile, err := os.Create("target_hashes.txt")
		if err == nil {
			for _, hashHex := range targetHashes {
				targetFile.WriteString(hashHex + "\n")
			}
			targetFile.Close()
		}

		fmt.Println("\nDump complete. You can use hashcat with these files:")
		fmt.Println("hashcat -m 300 target_hashes.txt candidates.txt")
		return
	}

	// CPU Mode
	fmt.Printf("CPU mode starting with %d threads...\n", *threads)
	bar := progressbar.Default(-1, "Hashing")

	var wg sync.WaitGroup
	chInfo := make(chan HashInfo, 100)

	for i := 0; i < *threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for phrase := range candidateCh {
				hashBytes := sha1.Sum([]byte(phrase))
				if targetHashesMap[hashBytes] {
					hashHex := hex.EncodeToString(hashBytes[:])
					chInfo <- HashInfo{Phrase: phrase, Hash: hashHex}
				}
				bar.Add(1)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(chInfo)
	}()

	for info := range chInfo {
		fmt.Printf("\nMatch found! Hash: %s Phrase: %s\n", info.Hash, info.Phrase)
		if err := SaveSeed(info.Phrase, info.Hash); err != nil {
			fmt.Printf("Error saving seed: %s\n", err)
		}
	}
	wgGen.Wait()
}