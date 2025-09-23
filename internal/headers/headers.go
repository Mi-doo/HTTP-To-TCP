package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() (h Headers) {
	return make(Headers)
}

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
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

	//Check for runes

	key := strings.ToLower(strings.Split(arrString[0], ":")[0])
	val := arrString[1]

	//Check if key already exits
	v, ok := h[key]
	if ok {
		val = strings.Join([]string{v, val}, ",")
	}

	h[key] = val

	return len(timedString) + len("\r\n"), false, nil
}
