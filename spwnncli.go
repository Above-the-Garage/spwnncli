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
	scoring      spwnn.PairScoring
	lengthMatch  spwnn.LengthMatch
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

func goCorrectSpelling(word string, noisy bool) {
	//start := time.Now()

	dict := getDict()
	correctedWords, _ := spwnn.CorrectSpelling(dict, word, lengthMatch, scoring)

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

// flagWasSet reports whether the named flag was given on the command line
func flagWasSet(name string) bool {
	set := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			set = true
		}
	})
	return set
}

func prompt() {
	fmt.Printf("\nCommand or word: ")
}

// A misspelling that should correct to one dictionary word.
type misspellingCase struct {
	typo     string
	expected string
}

// Typos whose only best match is the intended word.
// Dropped letters, extra letters, and swapped letters are all represented.
func misspellingCases() []misspellingCase {
	return []misspellingCase{
		{"ontemplation", "contemplation"},
		{"templation", "contemplation"},
		{"conteplation", "contemplation"},
		{"contemplationn", "contemplation"},
		{"seperate", "separate"},
		{"occassion", "occasion"},
		{"occurence", "occurrence"},
		{"accomodate", "accommodate"},
		{"necesary", "necessary"},
		{"enviroment", "environment"},
		{"goverment", "government"},
		{"tommorow", "tomorrow"},
		{"begining", "beginning"},
		{"calender", "calendar"},
		{"arguement", "argument"},
		{"experiance", "experience"},
		{"libary", "library"},
		{"mispell", "misspell"},
		{"writting", "writing"},
		{"suprise", "surprise"},
	}
}

// checkMisspellings reports typos that do not correct to exactly the intended word.
func checkMisspellings(dict *spwnn.SpwnnDictionary) int {
	cases := misspellingCases()
	failed := 0
	fmt.Printf("Misspellings that should correct to one word:\n")
	for _, c := range cases {
		want := "_" + c.expected + "_"
		results, _ := spwnn.CorrectSpelling(dict, c.typo, spwnn.AnyLength, scoring)
		if len(results) == 1 && results[0].Word == want {
			fmt.Printf("  ok    %s -> %s  %d%%\n", c.typo, results[0].Word, percentage(results[0].Score))
			continue
		}
		failed++
		fmt.Printf("  FAIL  %s expected %s\n", c.typo, want)
		printResults(c.typo, results)
	}
	fmt.Printf("Misspellings: %d checked, %d failed\n", len(cases), failed)
	return failed
}

func main() {

	wordToCorrect := flag.String("word", "", "a word to spelling correct")
	dictFlag := flag.String("dict", "knownWords.txt", "dictionary file to use")
	lengthFlag := flag.String("length", "any", "which winners to keep: any length, or the same length as the input")
	scoringFlag := flag.String("scoring", "capped", "how repeated letter pairs are scored: capped or original")
	testFlag := flag.Bool("test", false, "run self-test (each word should correct to itself)")
	flag.Parse()

	var err error
	scoring, err = parseScoring(*scoringFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	lengthMatch, err = parseLength(*lengthFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	// The self-test expects exactly one winner per word, so it has always
	// kept only same-length words; do that unless -length was given.
	if *testFlag && !flagWasSet("length") {
		lengthMatch = spwnn.SameLength
		*lengthFlag = "same"
	}
	fmt.Printf("Pair scoring: %s; length match: %s\n", *scoringFlag, *lengthFlag)

	dictFilename = *dictFlag
	dict := spwnn.ReadDictionary(dictFilename, true)

	if *testFlag {
		failed := checkMisspellings(dict)
		benchmarkParallel(dict.Words(), "", true)
		if failed != 0 {
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(*wordToCorrect) != 0 {
		correctedWords, _ := spwnn.CorrectSpelling(dict, *wordToCorrect, lengthMatch, scoring)
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
				correctedWords, wordsTouched := spwnn.CorrectSpelling(dict, input, lengthMatch, scoring)
				printResults(input, correctedWords)
				fmt.Printf("Words Touched = %d%%\n", percentage(float64(wordsTouched)/float64(dict.WordCount())))
			}
		}
		prompt()
	}
	os.Exit(0)
}

// parseScoring converts the -scoring flag text to a spwnn.PairScoring
func parseScoring(name string) (spwnn.PairScoring, error) {
	switch name {
	case "capped":
		return spwnn.CappedPairs, nil
	case "original":
		return spwnn.OriginalPairs, nil
	}
	return spwnn.CappedPairs, fmt.Errorf("unknown scoring %q (want capped or original)", name)
}

// parseLength converts the -length flag text to a spwnn.LengthMatch
func parseLength(name string) (spwnn.LengthMatch, error) {
	switch name {
	case "any":
		return spwnn.AnyLength, nil
	case "same":
		return spwnn.SameLength, nil
	}
	return spwnn.AnyLength, fmt.Errorf("unknown length %q (want any or same)", name)
}
