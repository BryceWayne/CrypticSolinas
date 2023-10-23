package main

import (
	"strconv"
)

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
