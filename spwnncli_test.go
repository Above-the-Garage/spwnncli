package main

import (
	"testing"

	"github.com/above-the-garage/spwnn"
	"github.com/stretchr/testify/assert"
)

func TestCorrector(t *testing.T) {
	InitDict()
	wordToTest := "stephen"
	correctedWords, touched := spwnn.CorrectSpelling(testDict, wordToTest)
	printResults(wordToTest, correctedWords)
	assert.NotEqual(t, 0, touched)
	assert.NotEqual(t, 0, len(correctedWords))
}
