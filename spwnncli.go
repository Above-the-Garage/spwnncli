/*
Package main - spwnncli
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/above-the-garage/spwnn"
)

func wordPresent(word string, results []spwnn.SpwnnResult) bool {
	for _, result := range results {
		if result.Word == word {
			return true
		}
	}
	return false
}

func percentage(val float64) int {
	return int(val * 100)
}

func printResults(word string, results []spwnn.SpwnnResult) {
	for _, result := range results {
		fmt.Printf("  %d%%\t%s\n", percentage(result.Score), result.Word)
	}
}

// Validate checks that all words correct to themselves - serial version
// Prints those that don't
func validate(dict *spwnn.SpwnnDictionary, letters string, noisy bool) {
	letters = spwnn.RemoveSpaces(letters)
	if len(letters) == 0 {
		letters = "abcdefghijklmnopqrstuvwxyz"
	}
	fmt.Printf("Letters = '%s'; noisy = %v\n", letters, noisy)
	for _, word := range spwnn.GetWords(dict) {
		testLetterIndex := 0
		if word[0] == '_' {
			testLetterIndex++
		}
		testLetter := fmt.Sprintf("%c", word[testLetterIndex])
		if strings.ContainsAny(testLetter, letters) {
			correctedWords, _ := spwnn.CorrectSpelling(dict, word)
			if noisy && !wordPresent(word, correctedWords) {
				fmt.Printf("Validation:  '%s' miscorrected to '%v'\n", word, correctedWords)
			}
		}
	}
}

// checkForTies checks that all words correct to themselves - serial version
// Counts number of errors
func checkForTies(dict *spwnn.SpwnnDictionary, letters string, noisy bool) {
	letters = spwnn.RemoveSpaces(letters)
	if len(letters) == 0 {
		letters = "abcdefghijklmnopqrstuvwxyz"
	}
	fmt.Printf("Letters = '%s'; noisy = %v\n", letters, noisy)
	totalWordsTied := 0
	totalTies := 0
	for _, word := range spwnn.GetWords(dict) {
		testLetterIndex := 0
		if word[0] == '_' {
			testLetterIndex++
		}
		testLetter := fmt.Sprintf("%c", word[testLetterIndex])
		if strings.ContainsAny(testLetter, letters) {
			correctedWords, _ := spwnn.CorrectSpelling(dict, word)
			numResults := len(correctedWords)
			if numResults > 1 {
				totalWordsTied++
				totalTies = totalTies + numResults
				if noisy {
					fmt.Printf("Validation:  '%s' tied (%d): '%v'\n", word, numResults, correctedWords)
				}
			}
		}
	}
	fmt.Printf("Total Words Tied: %d; Total Ties: %d\n", totalWordsTied, totalTies)
}

// Benchmark prints the time to Validate words that start with the letters provided
func benchmark(dict *spwnn.SpwnnDictionary, input string) {
	letters := spwnn.RemoveSpaces(input)
	if len(letters) == 0 {
		letters = "abcdefghijklmnopqrstuvwxyz"
	}
	start := time.Now()
	validate(dict, letters, false)
	elapsed := time.Since(start)
	fmt.Printf("Elapsed time %s\n", elapsed)
}

//
// Parallel version of benchmarking code
//

var (
	mu    sync.Mutex
	dicts []*spwnn.SpwnnDictionary
)

// Manage a set of dictionaries, growing the list on demand
func getDict() *spwnn.SpwnnDictionary {
	mu.Lock()
	defer mu.Unlock()
	for i, dict := range dicts {
		if dict != nil {
			res := dict
			// Unlink to make it inaccessible
			// Go routine will add it back when done
			dicts[i] = nil
			return res
		}
	}
	newDict := spwnn.ReadDictionary(false)
	return newDict
}

func releaseDict(dict *spwnn.SpwnnDictionary) {
	mu.Lock()
	defer mu.Unlock()
	for i, d := range dicts {
		if d == nil {
			dicts[i] = dict
			return
		}
	}
	// All full, must have created a new dictionary, save for reuse
	dicts = append(dicts, dict)
}

var tokens = make(chan int, runtime.GOMAXPROCS(0))

func goCorrectSpelling(word string, noisy bool) {
	start := time.Now()

	dict := getDict()
	correctedWords, _ := spwnn.CorrectSpelling(dict, word)

	if noisy {
		fmt.Printf("%s: %s\n", word, time.Since(start))
	}
	if !wordPresent(word, correctedWords) {
		fmt.Printf("Validation:  '%s' miscorrected to '%v'\n", word, correctedWords)
	}

	// Let next go routine run
	releaseDict(dict)
	<-tokens
}

func benchmarkParallel(words []string, input string, noisy bool) {
	letters := spwnn.RemoveSpaces(input)
	if len(letters) == 0 {
		letters = "abcdefghijklmnopqrstuvwxyz"
	}
	fmt.Printf("Letters = '%s'; noisy = %v\n", letters, noisy)

	start := time.Now()

	for _, word := range words {
		testLetterIndex := 0
		if word[0] == '_' {
			testLetterIndex++
		}
		testLetter := fmt.Sprintf("%c", word[testLetterIndex])
		if strings.ContainsAny(testLetter, letters) {
			tokens <- len(tokens)
			go goCorrectSpelling(word, noisy)
		}
	}

	// Wait for all go routines to complete (this may be considered bad form!)
	for len(tokens) > 0 {
		time.Sleep(10 * time.Millisecond)
	}
	elapsed := time.Since(start)

	fmt.Printf("Elapsed time %s; GOMAXPROCS %d\n", elapsed, runtime.GOMAXPROCS(0))
}

func handleCommand(dict *spwnn.SpwnnDictionary, input string) {
	if len(input) == 0 {
		fmt.Printf("Say what?\n")
		return
	}

	cmd := input[0]
	input = input[1:]

	switch {

	case cmd == 'a':
		spwnn.SetAccuracy(dict, input)

	case cmd == 'b':
		benchmark(dict, input)

	case cmd == 'e', cmd == 'q':
		fmt.Printf("Bye!\n")
		os.Exit(0)

	case cmd == 'g':
		benchmarkParallel(spwnn.GetWords(dict), input, false)

	case cmd == 'm':
		maxSize := spwnn.MaxIndexSize(dict)
		fmt.Printf("max index = %d\n", maxSize)

	case cmd == 'p':
		spwnn.PrintNeuron(dict, input)

	case cmd == 's':
		spwnn.PrintIndexSizes(dict)

	case cmd == 't':
		input = spwnn.RemoveSpaces(input)
		checkForTies(dict, input, true)

	case cmd == 'v':
		input = spwnn.RemoveSpaces(input)
		validate(dict, input, true)

	default:
		fmt.Printf("say what?\n")
	}
}

func prompt() {
	fmt.Printf("\nCommand or word: ")
}

var testDict *spwnn.SpwnnDictionary

// InitDict reads the dictionary; it is exposed here
// so unit tests can init
func InitDict() {
	testDict = spwnn.ReadDictionary(true)
}

func main() {

	dict := spwnn.ReadDictionary(true)

	wordToCorrect := flag.String("word", "", "a word to spelling correct")
	flag.Parse()

	if len(*wordToCorrect) != 0 {
		correctedWords, _ := spwnn.CorrectSpelling(dict, *wordToCorrect)
		printResults(*wordToCorrect, correctedWords)
		os.Exit(0)
	}

	scanner := bufio.NewScanner(os.Stdin)
	prompt()
	for scanner.Scan() {
		var input string
		input = scanner.Text()
		if input[0] == '-' {
			handleCommand(dict, input[1:])
		} else {
			correctedWords, wordsTouched := spwnn.CorrectSpelling(dict, input)
			printResults(input, correctedWords)
			fmt.Printf("Words Touched = %d%%\n", percentage(float64(wordsTouched)/float64(spwnn.GetWordCount(dict))))
		}
		prompt()
	}
	os.Exit(0)
}
