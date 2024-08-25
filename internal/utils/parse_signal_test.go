package utils_test

import (
	"strconv"
	"syscall"
	"testing"

	"github.com/toudi/joker/internal/utils"
)

func TestParseSignalFromString(t *testing.T) {
	type testCase struct {
		input    string
		expected syscall.Signal
	}

	for _, test := range []testCase{
		{
			input:    "TERM",
			expected: syscall.SIGTERM,
		},
		{
			input:    "SIGTERM",
			expected: syscall.SIGTERM,
		},
		{
			input:    "QUIT",
			expected: syscall.SIGQUIT,
		},
		{
			input:    "SIGQUIT",
			expected: syscall.SIGQUIT,
		},
		{
			input:    "XYZ",
			expected: syscall.SIGTERM,
		},
	} {
		test := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()

			if parsedSignal := utils.ParseSignalFromString(test.input); parsedSignal != test.expected {
				t.Errorf("got %v; expected %v\n", parsedSignal, test.expected)
			}
		})
	}
}

func TestParseSignalFromInt(t *testing.T) {
	type testCase struct {
		input    int
		expected syscall.Signal
	}

	for _, test := range []testCase{
		{
			input:    15,
			expected: syscall.SIGTERM,
		},
		{
			input:    2,
			expected: syscall.SIGINT,
		},
		{
			input:    9,
			expected: syscall.SIGKILL,
		},
		{
			input:    -5,
			expected: syscall.SIGTERM,
		},
	} {
		test := test
		t.Run(strconv.Itoa(test.input), func(t *testing.T) {
			t.Parallel()

			if parsedSignal := utils.ParseSignalFromInt(test.input); parsedSignal != test.expected {
				t.Errorf("got %v; expected %v\n", parsedSignal, test.expected)
			}
		})
	}
}
