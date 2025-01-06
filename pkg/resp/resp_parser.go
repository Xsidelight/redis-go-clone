package resp

import (
	"errors"
	"strconv"
	"strings"
)

func DeserializeRESP(input string) (any, error) {
	value, _, err := parseOneValue(input)
	return value, err
}

// parseOneValue reads a single RESP value from `input` and returns:
// - the parsed Go value
// - the remaining unconsumed part of `input`
// - any error
func parseOneValue(input string) (any, string, error) {
	// Ensure we have something
	if len(input) == 0 {
		return nil, "", errors.New("empty input")
	}

	switch input[0] {
	case '+': // Simple string
		return parseSimpleStringValue(input)
	case '-': // Error
		return parseErrorValue(input)
	case ':': // Integer
		return parseIntegerValue(input)
	case '$': // Bulk string
		return parseBulkStringValue(input)
	case '*': // Array
		return parseArrayValue(input)
	default:
		return nil, "", errors.New("unsupported RESP type")
	}
}

// parseSimpleStringValue returns (string, leftover, error).
func parseSimpleStringValue(input string) (any, string, error) {
	endIdx := strings.Index(input, "\r\n")
	if endIdx == -1 {
		return nil, "", errors.New("invalid RESP format for simple string")
	}
	// e.g. +OK\r\n => "OK"
	value := input[1:endIdx]
	remaining := input[endIdx+2:]
	return value, remaining, nil
}

// parseErrorValue is similar to parseSimpleStringValue but returns the error message (string).
func parseErrorValue(input string) (any, string, error) {
	endIdx := strings.Index(input, "\r\n")
	if endIdx == -1 {
		return nil, "", errors.New("invalid RESP format for error")
	}
	value := input[1:endIdx]
	remaining := input[endIdx+2:]
	return value, remaining, nil
}

// parseIntegerValue parses RESP integer like `:42\r\n`.
func parseIntegerValue(input string) (any, string, error) {
	endIdx := strings.Index(input, "\r\n")
	if endIdx == -1 {
		return nil, "", errors.New("invalid RESP format for integer")
	}

	valueStr := input[1:endIdx]
	remaining := input[endIdx+2:]

	v, err := strconv.Atoi(valueStr)
	if err != nil {
		return nil, "", errors.New("invalid integer value")
	}
	return v, remaining, nil
}

// parseBulkStringValue parses `$<length>\r\n<content>\r\n`
func parseBulkStringValue(input string) (any, string, error) {
	// 1) find the first line up to \r\n
	firstLineEnd := strings.Index(input, "\r\n")
	if firstLineEnd == -1 {
		return nil, "", errors.New("invalid RESP format for bulk string length")
	}

	// "$4\r\n1234\r\n..."
	lengthStr := input[1:firstLineEnd] // the part after '$' but before "\r\n"
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, "", errors.New("invalid bulk string length")
	}

	remaining := input[firstLineEnd+2:] // move past "$4\r\n"

	if length == -1 {
		// $-1\r\n => nil
		return nil, remaining, nil
	}

	// We must have at least `length` characters + 2 bytes for the trailing "\r\n"
	if len(remaining) < length+2 {
		return nil, "", errors.New("not enough data for bulk string")
	}

	content := remaining[:length]
	trailer := remaining[length : length+2]
	if trailer != "\r\n" {
		return nil, "", errors.New("invalid bulk string trailer")
	}

	newRemaining := remaining[length+2:]

	// Attempt to parse content as int
	if intVal, err := strconv.Atoi(content); err == nil {
		return intVal, newRemaining, nil
	}

	// Else store as string
	return content, newRemaining, nil
}

// parseArrayValue parses `*3\r\n...`
func parseArrayValue(input string) (any, string, error) {
	// 1) find the first line up to \r\n
	firstLineEnd := strings.Index(input, "\r\n")
	if firstLineEnd == -1 {
		return nil, "", errors.New("invalid RESP format for array")
	}

	// "*3\r\n..."
	countStr := input[1:firstLineEnd] // the part after '*' but before "\r\n"
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return nil, "", errors.New("invalid array count")
	}

	remaining := input[firstLineEnd+2:] // move past "*3\r\n"

	if count == -1 {
		// *-1 => nil array
		return nil, remaining, nil
	}

	elements := make([]any, 0, count)

	// Parse each element
	for i := 0; i < count; i++ {
		val, newRem, err := parseOneValue(remaining)
		if err != nil {
			return nil, "", err
		}
		elements = append(elements, val)
		remaining = newRem
	}

	return elements, remaining, nil
}

func SerializeRESP(data any, isGet bool) string {
	switch v := data.(type) {
	case string:
		if isGet {
			return serializeBulkString(v)
		}

		return serializeSimpleString(v)
	case int:
		return serializeInteger(v)
	case []any:
		return serializeArray(v)
	case error:
		return serializeError(v.Error())
	case []byte:
		if len(v) == 0 {
			return serializeNull()
		}
		return serializeBulkString(string(v))
	case nil:
		return serializeNull()
	default:
		return serializeError("unsupported RESP type")
	}
}

func serializeSimpleString(value string) string {
	return "+" + value + "\r\n"
}

func serializeInteger(value int) string {
	return ":" + strconv.Itoa(value) + "\r\n"
}

func serializeError(message string) string {
	return "-" + message + "\r\n"
}

func serializeNull() string {
	return "$-1\r\n"
}

func serializeBulkString(value string) string {
	if value == "" {
		return serializeNull()
	}
	return "$" + strconv.Itoa(len(value)) + "\r\n" + value + "\r\n"
}

func serializeArray(array []any) string {
	if len(array) == 0 {
		return "*0\r\n"
	}

	var sb strings.Builder
	sb.WriteString("*")
	sb.WriteString(strconv.Itoa(len(array)))
	sb.WriteString("\r\n")

	for _, element := range array {
		sb.WriteString(SerializeRESP(element, false))
	}
	return sb.String()
}
