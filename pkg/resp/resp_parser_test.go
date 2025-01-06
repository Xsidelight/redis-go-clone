package resp

import (
	"errors"
	"reflect"
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
		{input: "$11\r\nHello World\r\n", expected: "Hello World", hasError: false},
		{input: "$-1\r\n", expected: nil, hasError: false},
		{input: "+hello world\r\n", expected: "hello world", hasError: false},
		{input: "*3\r\n$3\r\nset\r\n$5\r\ntestv\r\n$4\r\n1234\r\n", expected: []any{"set", "testv", 1234}, hasError: false},
		{input: "+NoEndLine", expected: nil, hasError: true},
		{input: "", expected: nil, hasError: true},
		{input: ":1000\r\n", expected: 1000, hasError: false},
		{input: "*-1\r\n", expected: nil, hasError: false},
		{input: "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", expected: []any{"foo", "bar"}, hasError: false},
		{input: "*2\r\n:1\r\n:2\r\n", expected: []any{1, 2}, hasError: false},
		{input: "*2\r\n$3\r\nfoo\r\n:42\r\n", expected: []any{"foo", 42}, hasError: false},
		{input: "*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n", expected: []any{"foo", nil, "bar"}, hasError: false},
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
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v for input %q, got %v", tc.expected, tc.input, result)
			}
		}
	}
}

func TestSerializeRESP(t *testing.T) {
	testCases := []struct {
		input    any
		expected string
	}{
		// Valid cases
		{input: "OK", expected: "+OK\r\n"},
		{input: 42, expected: ":42\r\n"},
		{input: error(errors.New("Error message")), expected: "-Error message\r\n"},
		{input: []byte("Hello World"), expected: "$11\r\nHello World\r\n"},
		{input: nil, expected: "$-1\r\n"},
		{input: []any{"SET", "key", "value"}, expected: "*3\r\n+SET\r\n+key\r\n+value\r\n"},
		{input: []any{1, "two", nil}, expected: "*3\r\n:1\r\n+two\r\n$-1\r\n"},
		{input: []byte(""), expected: "$-1\r\n"},
		{input: []any{}, expected: "*0\r\n"},                                                     // Empty array
		{input: []any{"foo", []any{"bar", 42}}, expected: "*2\r\n+foo\r\n*2\r\n+bar\r\n:42\r\n"}, // Nested array
		{input: "hello", expected: "+hello\r\n"},
		{input: -1, expected: ":-1\r\n"},
		{input: []any{"foo", nil, "bar"}, expected: "*3\r\n+foo\r\n$-1\r\n+bar\r\n"},

		// Invalid cases
		{input: struct{}{}, expected: "-unsupported RESP type\r\n"}, // Unsupported type
	}

	for _, tc := range testCases {
		result := SerializeRESP(tc.input, false)
		if result != tc.expected {
			t.Errorf("SerializeRESP(%v) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}
