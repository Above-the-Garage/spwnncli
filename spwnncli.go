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
		fmt.Printf("  %d%%\t%d\t%s\n", percentage(result.Score), result.LenDiff, result.Word)
	}
}

//
// Parallel version of benchmarking code
//

var (
	mu           sync.Mutex
	dicts        []*spwnn.SpwnnDictionary
	dictFilename string
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
	newDict := spwnn.ReadDictionary(dictFilename, false)
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

func goCorrectSpelling(word string, noisy bool, strictLen bool) {
	//start := time.Now()

	dict := getDict()
	correctedWords, _ := spwnn.CorrectSpelling(dict, word, strictLen)

	if len(correctedWords) != 1 {
		fmt.Printf("Parallel validation:  '%s' could be '%v'\n", word, correctedWords)
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
			go goCorrectSpelling(word, noisy, true /* strictLen */)
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
	input = strings.TrimSpace(input[1:])

	switch {

	case cmd == 'e', cmd == 'q':
		fmt.Printf("Bye!\n")
		os.Exit(0)

	case cmd == 'g':
		benchmarkParallel(dict.Words(), input, true)

	case cmd == 'm':
		maxSize := spwnn.MaxIndexSize(dict)
		fmt.Printf("max index = %d\n", maxSize)

	case cmd == 'p':
		spwnn.PrintNeuron(dict, input)

	case cmd == 's':
		spwnn.PrintIndexSizes(dict)

	default:
		fmt.Printf("say what?\n")
	}
}

func prompt() {
	fmt.Printf("\nCommand or word: ")
}

func main() {

	wordToCorrect := flag.String("word", "", "a word to spelling correct")
	dictFlag := flag.String("dict", "knownWords.txt", "dictionary file to use")
	flag.Parse()

	dictFilename = *dictFlag
	dict := spwnn.ReadDictionary(dictFilename, true)

	if len(*wordToCorrect) != 0 {
		correctedWords, _ := spwnn.CorrectSpelling(dict, *wordToCorrect, false /* strictLen */)
		printResults(*wordToCorrect, correctedWords)
		os.Exit(0)
	}

	scanner := bufio.NewScanner(os.Stdin)
	prompt()
	for scanner.Scan() {
		var input string
		input = scanner.Text()
		if len(input) != 0 {
			if input[0] == '-' {
				handleCommand(dict, input[1:])
			} else {
				correctedWords, wordsTouched := spwnn.CorrectSpelling(dict, input, false /* strictLen */)
				printResults(input, correctedWords)
				fmt.Printf("Words Touched = %d%%\n", percentage(float64(wordsTouched)/float64(dict.WordCount())))
			}
		}
		prompt()
	}
	os.Exit(0)
}
