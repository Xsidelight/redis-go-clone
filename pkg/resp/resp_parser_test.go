package resp

import (
	"testing"
)

type testCase struct {
	input    string
	expected any
	hasError bool
}

func TestDeserializeRESP(t *testing.T) {
	testCases := []testCase{
		{input: "+OK\r\n", expected: "OK", hasError: false},
		{input: "-Error message\r\n", expected: "Error message", hasError: false},
		{input: "$0\r\n\r\n", expected: "", hasError: false},
		{input: "$-1\r\n", expected: nil, hasError: false},
		{input: "+hello world\r\n", expected: "hello world", hasError: false},
		{input: "+NoEndLine", expected: nil, hasError: true},
		{input: "", expected: nil, hasError: true},
	}

	for _, tc := range testCases {
		result, err := DeserializeRESP(tc.input)

		if tc.hasError {
			if err == nil {
				t.Errorf("expected error for input %q, got nil", tc.input)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tc.input, err)
			}
			if result != tc.expected {
				t.Errorf("expected %v for input %q, got %v", tc.expected, tc.input, result)
			}
		}
	}
}
