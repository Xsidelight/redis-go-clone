package resp

import (
	"errors"
	"strconv"
	"strings"
)

func DeserializeRESP(s string) (any, error) {
	if strings.HasPrefix(s, "+") {
		return parseSimpleString(s)
	} else if strings.HasPrefix(s, "-") {
		return parseError(s)
	} else if strings.HasPrefix(s, ":") {
		return parseSimpleInteger(s)
	} else if strings.HasPrefix(s, "$") {
		return parseBulkString(s)
	} else if strings.HasPrefix(s, "*") {
		return parseArray(s)
	}
	return nil, errors.New("unsupported RESP type")
}

func parseSimpleString(s string) (string, error) {
	if !strings.HasSuffix(s, "\r\n") {
		return "", errors.New("invalid RESP format for Simple String")
	}
	return strings.TrimSuffix(s[1:], "\r\n"), nil
}

func parseError(s string) (string, error) {
	if !strings.HasSuffix(s, "\r\n") {
		return "", errors.New("invalid RESP format for Error")
	}
	return strings.TrimSuffix(s[1:], "\r\n"), nil
}

func parseSimpleInteger(s string) (int, error) {
	if !strings.HasSuffix(s, "\r\n") {
		return 0, errors.New("invalid RESP format for Error")
	}

	s = strings.TrimSuffix(s[1:], "\r\n")
	r, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("error parsing to integer")
	}
	return r, nil
}

func parseBulkString(s string) (any, error) {
	if !strings.HasPrefix(s, "$") {
		return nil, errors.New("invalid RESP Bulk String format")
	}

	endOfLength := strings.Index(s, "\r\n")
	if endOfLength == -1 {
		return nil, errors.New("invalid RESP Bulk String format")
	}

	lengthStr := s[1:endOfLength]
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, errors.New("invalid Bulk String length")
	}

	if length == -1 {
		return nil, nil
	}

	startOfContent := endOfLength + 2
	endOfContent := startOfContent + length
	if endOfContent+2 > len(s) || !strings.HasSuffix(s[endOfContent:endOfContent+2], "\r\n") {
		return nil, errors.New("invalid RESP Bulk String content")
	}

	return s[startOfContent:endOfContent], nil
}

func parseArray(s string) ([]any, error) {
	if !strings.HasPrefix(s, "*") {
		return nil, errors.New("invalid RESP Array format")
	}

	endOfCount := strings.Index(s, "\r\n")
	if endOfCount == -1 {
		return nil, errors.New("invalid RESP Array format")
	}

	countStr := s[1:endOfCount]
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return nil, errors.New("invalid Array count")
	}

	if count == -1 {
		return nil, nil
	}

	elements := make([]any, 0, count)
	remaining := s[endOfCount+2:]
	for i := 0; i < count; i++ {
		element, err := DeserializeRESP(remaining)
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)

		// Adjust remaining to skip the parsed element
		pos := strings.Index(remaining, "\r\n") + 2
		remaining = remaining[pos:]
	}

	return elements, nil
}

func SerializeRESP(data any) string {
	switch v := data.(type) {
	case string:
		return serializeSimpleString(v)
	case int:
		return serializeSimpleInteger(v)
	case []any:
		return serializeArray(v)
	case error:
		return serializeError(v.Error())
	case []byte:
		return serializeBulkString(string(v))
	case nil:
		return serializeNull()
	default:
		return serializeError("unsupported RESP type")
	}
}

func serializeSimpleString(v string) string {
	return "+" + v + "\r\n"
}

func serializeSimpleInteger(v int) string {
	return ":" + strconv.Itoa(v) + "\r\n"
}

func serializeError(s string) string {
	return "-" + s + "\r\n"
}

func serializeNull() string {
	return "$-1\r\n"
}

func serializeBulkString(s string) string {
	if s == "" {
		return serializeNull()
	}
	return "$" + strconv.Itoa(len(s)) + "\r\n" + s + "\r\n"
}

func serializeArray(a []any) string {
	var sb strings.Builder
	sb.WriteString("*")
	sb.WriteString(strconv.Itoa(len(a)))
	sb.WriteString("\r\n")

	for _, element := range a {
		sb.WriteString(SerializeRESP(element))
	}

	return sb.String()
}
