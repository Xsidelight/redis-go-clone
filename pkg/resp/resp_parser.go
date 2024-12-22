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

	var elements []any
	remaining := s[endOfCount+2:]

	for i := 0; i < count; i++ {
		element, err := DeserializeRESP(remaining)
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)

		nextStart := strings.Index(remaining, "\r\n") + 2
		if nextStart == 1 {
			return nil, errors.New("malformed RESP Array")
		}
		remaining = remaining[nextStart:]
	}

	return elements, nil
}

func SerializeRESP() {
	// TODO: Implement serialization logic
}
