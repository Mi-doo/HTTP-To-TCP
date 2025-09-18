package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	request := strings.Split(string(content), "\r\n")
	startLine := request[0]

	line, err := parseRequestLine(startLine)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *line,
	}, nil
}

func parseRequestLine(content string) (*RequestLine, error) {

	if !strings.Contains(content, "HTTP/1.1") {
		msg := "We only support HTTP/1.1 for now"
		return nil, fmt.Errorf("%s", msg)
	}

	r := strings.Split(content, " ")
	if len(r) != 3 {
		msg := "Request line is not correct"
		return nil, fmt.Errorf("%s", msg)
	}

	version := strings.Split(r[2], "/")[1]

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: r[1],
		Method:        r[0],
	}, nil
}
