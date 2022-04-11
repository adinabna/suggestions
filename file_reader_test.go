package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFileWithOneEntry(t *testing.T) {
	data := readFile("input/one_entry.csv")
	assert.NotNil(t, data)
	assert.Equal(t, len(data), 1)
	assert.Equal(t, len(data[0]), 3)
}

func TestReadFileWithMultipleEntries(t *testing.T) {
	data := readFile("input/example.csv")
	assert.NotNil(t, data)
	assert.Equal(t, len(data), 2)
}