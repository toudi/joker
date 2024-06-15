package joker

import (
	"fmt"

	"github.com/mgutz/ansi"
)

type StreamLine struct {
	Service string
	Stderr  bool
	Line    string
}

var reset = ansi.ColorCode("reset")

func (j *Joker) StreamHandler() {
	go func() {
		for entry := range j.streamChan {
			var streamType = "stdout"
			if entry.Stderr {
				streamType = "stderr"
			}
			fmt.Printf(
				"%s%s | %s | %s%s\n",
				getColorFunc(entry.Service),
				entry.Service,
				streamType,
				entry.Line,
				reset,
			)
		}
	}()
}

var colorForService map[string]string
var availableStyles = []string{
	"red",
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
	"white",
}
var lastColorUsed int = -1

func getColorFunc(service string) string {
	if colorForService == nil {
		colorForService = make(map[string]string)
	}

	if _, exists := colorForService[service]; !exists {
		lastColorUsed += 1
		colorForService[service] = ansi.ColorCode(
			availableStyles[lastColorUsed%len(availableStyles)],
		)
	}

	return colorForService[service]
}
