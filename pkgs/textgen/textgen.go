package textgen

import (
	"encoding/json"
	"math/rand"
	"os"
	"strings"
)

func LoadModel(filename string) (map[string][]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var markovChain map[string][]string
	if err := json.Unmarshal(data, &markovChain); err != nil {
		return nil, err
	}

	parsedChain := make(map[string][]string)
	for key, value := range markovChain {
		cleanKey := strings.Trim(key, "()'',")
		parsedChain[cleanKey] = value
	}

	return parsedChain, nil
}

func GenerateSentence(markovChain map[string][]string, length int, seed string) string {
	var seedWords []string
	if seed == "" {
		keys := make([]string, 0, len(markovChain))
		for key := range markovChain {
			keys = append(keys, key)
		}
		seedWords = strings.Fields(keys[rand.Intn(len(keys))])
	} else {
		seedWords = strings.Fields(seed)
	}

	sentence := seedWords
	for i := 0; i < length-len(seedWords); i++ {
		lastWords := strings.Join(sentence[len(sentence)-len(seedWords):], " ")
		nextWords, exists := markovChain[lastWords]
		if !exists || len(nextWords) == 0 {
			break
		}
		nextWord := nextWords[rand.Intn(len(nextWords))]
		sentence = append(sentence, nextWord)
	}

	return strings.Join(sentence, " ")
}
