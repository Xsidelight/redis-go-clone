package redis_command

import (
	"errors"
	"fmt"
	"redis-go-clone/internal/model"
	"sync"
	"time"
)

func Set(cmdArray []any, storedData map[string]model.StoredData, mu *sync.RWMutex) string {

	if len(cmdArray) < 3 {
		return "-ERR wrong number of arguments for SET\r\n"
	}

	key, ok := cmdArray[1].(string)
	if !ok {
		return "-ERR invalid argument for SET\r\n"
	}

	value := cmdArray[2]
	expiryTimestamp := int64(0)

	if len(cmdArray) > 3 {
		var err error
		expiryTimestamp, err = parseExpiryOptions(cmdArray[3:])
		if err != nil {
			return fmt.Sprintf("-ERR %s\r\n", err.Error())
		}
	}

	mu.Lock()
	defer mu.Unlock()
	storedData[key] = model.StoredData{Value: value, ExpiryDate: expiryTimestamp}

	return "+OK\r\n"
}

func parseExpiryOptions(options []any) (int64, error) {
	if len(options) < 2 {
		return 0, errors.New("invalid expiry arguments")
	}

	optionType, ok := options[0].(string)
	if !ok {
		return 0, errors.New("invalid expiry option type")
	}

	expiryValue, ok := options[1].(int64)
	if !ok {
		return 0, errors.New("invalid expiry value")
	}

	switch optionType {
	case "EX":
		return time.Now().Add(time.Duration(expiryValue) * time.Second).Unix(), nil
	case "PX":
		return time.Now().Add(time.Duration(expiryValue) * time.Millisecond).UnixMilli(), nil
	case "EXAT":
		if expiryValue <= 0 {
			return 0, errors.New("invalid Unix time for EXAT")
		}
		return expiryValue, nil
	case "PXAT":
		if expiryValue <= 0 {
			return 0, errors.New("invalid Unix time for PXAT")
		}
		return time.Unix(0, expiryValue*int64(time.Millisecond)).Unix(), nil
	default:
		return 0, errors.New("unsupported expiry option")
	}

}
