package resp

import (
	"errors"
	"strconv"
	"strings"
)

// DeserializeRESP parses a RESP-encoded string and returns the corresponding Go value.
func DeserializeRESP(input string) (any, error) {
	switch {
	case strings.HasPrefix(input, "+"):
		return parseSimpleString(input)
	case strings.HasPrefix(input, "-"):
		return parseError(input)
	case strings.HasPrefix(input, ":"):
		return parseInteger(input)
	case strings.HasPrefix(input, "$"):
		return parseBulkString(input)
	case strings.HasPrefix(input, "*"):
		return parseArray(input)
	default:
		return nil, errors.New("unsupported RESP type")
	}
}

// parseSimpleString parses a RESP Simple String.
func parseSimpleString(input string) (string, error) {
	if !strings.HasSuffix(input, "\r\n") {
		return "", errors.New("invalid RESP format for simple string")
	}
	return strings.TrimSuffix(input[1:], "\r\n"), nil
}

// parseError parses a RESP Error.
func parseError(input string) (string, error) {
	if !strings.HasSuffix(input, "\r\n") {
		return "", errors.New("invalid RESP format for error")
	}
	return strings.TrimSuffix(input[1:], "\r\n"), nil
}

// parseInteger parses a RESP Integer.
func parseInteger(input string) (int, error) {
	if !strings.HasSuffix(input, "\r\n") {
		return 0, errors.New("invalid RESP format for integer")
	}
	valueStr := strings.TrimSuffix(input[1:], "\r\n")
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, errors.New("invalid integer value")
	}
	return value, nil
}

// parseBulkString parses a RESP Bulk String.
func parseBulkString(input string) (any, error) {
	parts := strings.SplitN(input, "\r\n", 2)
	if len(parts) < 2 {
		return nil, errors.New("invalid RESP format for bulk string")
	}

	length, err := strconv.Atoi(parts[0][1:])
	if err != nil {
		return nil, errors.New("invalid bulk string length")
	}

	if length == -1 {
		return nil, nil
	}

	if len(parts[1]) < length+2 || !strings.HasSuffix(parts[1][:length+2], "\r\n") {
		return nil, errors.New("invalid bulk string content")
	}

	return parts[1][:length], nil
}

// parseArray parses a RESP Array.
func parseArray(input string) ([]any, error) {
	header := strings.SplitN(input, "\r\n", 2)
	if len(header) < 2 {
		return nil, errors.New("invalid RESP format for array")
	}

	count, err := strconv.Atoi(header[0][1:])
	if err != nil || count < -1 {
		return nil, errors.New("invalid array count")
	}

	if count == -1 {
		return nil, nil
	}

	elements := make([]any, 0, count)
	remaining := header[1]

	for i := 0; i < count; i++ {
		element, err := DeserializeRESP(remaining)
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)

		// Adjust the remaining input based on parsed element
		switch v := element.(type) {
		case string:
			remaining = strings.TrimPrefix(remaining, serializeBulkString(v))
		case []any:
			remaining = strings.TrimPrefix(remaining, serializeArray(v))
		case int:
			remaining = strings.TrimPrefix(remaining, serializeInteger(v))
		default:
			return nil, errors.New("unsupported element type in array")
		}
	}

	return elements, nil
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
