package core

import (
	"log"
	"plugin"
)

type Task struct{
	File		string
	Status		int
	Type		int
	Id			int
}

type KeyValue struct{
	Key			string
	Value		string
}

func ReadMapReduceFuncs(filename string) (func(string, string) []KeyValue, func(string, []string) string) {
	p, err := plugin.Open(filename)
	if err != nil {
		log.Fatalf("Can't open Plugin\n")
	}
	xmapf, err := p.Lookup("Map")
	if err != nil {
		log.Fatalf("Map not found\n")
	}
	mapf := xmapf.(func(string, string) []KeyValue)
	xreducef, err := p.Lookup("Reduce")
	if err != nil {
		log.Fatalf("Reduce not found\n")
	}
	reducef := xreducef.(func(string, []string) string)

	return mapf, reducef
} 