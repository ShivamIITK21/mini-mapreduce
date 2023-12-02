package main

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/ShivamIITK21/mini-mapreduce/core"
)

func Map(filename string, contents string) []core.KeyValue {
	// function to detect word separators.
	ff := func(r rune) bool { return !unicode.IsLetter(r) }

	// split contents into an array of words.
	words := strings.FieldsFunc(contents, ff)

	kva := []core.KeyValue{}
	for _, w := range words {
			kv := core.KeyValue{Key: w, Value: "1"}
			kva = append(kva, kv)
	}
	return kva
}


func Reduce(key string, values []string) string {
	// return the number of occurrences of this word.
	return strconv.Itoa(len(values))
}
