package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() (h Headers) {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if !strings.Contains(string(data), "\r\n") {
		return 0, false, nil
	}

	if strings.HasPrefix(string(data), "\r\n") {
		return 0, true, nil
	}

	splitString := strings.Split(string(data), "\r\n")
	timedString := strings.TrimSpace(splitString[0])
	arrString := strings.Fields(timedString)
	if len(arrString) > 2 {
		msg := "Not a valid header"
		err := fmt.Errorf("%s", msg)
		return 0, false, err
	}

	key := strings.Split(arrString[0], ":")[0]
	val := arrString[1]
	h[key] = val

	return len(timedString) + len("\r\n"), false, nil
}
