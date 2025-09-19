package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

const bufferSize = 8

type state string

const (
	initialized state = "initialized"
	done              = "done"
)

type Request struct {
	RequestLine RequestLine
	state       state
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type ChunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *ChunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if !strings.Contains(string(data), "\r\n") {
		return 0, nil
	}

	startLine := strings.Split(string(data), "\r\n")[0]
	line, err := parseRequestLine(startLine)
	if err != nil {
		return 0, err
	}

	r.RequestLine = *line

	return len(data), nil
}

func RequestFromReader(reader *ChunkReader) (*Request, error) {

	buff := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	r := &Request{
		state: "initialized",
	}
	for r.state != "done" {
		if len(buff) == cap(buff) {
			buff2 := make([]byte, len(buff)*2, cap(buff)*2)
			copy(buff2, buff)
			buff = buff2
		}

		n, err := reader.Read(buff[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.state = "done"
				break
			}
			return nil, err
		}

		readToIndex += n

		numBytesParsed, err := r.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}

		s := make([]byte, numBytesParsed, numBytesParsed)
		m := copy(buff[:numBytesParsed], s)

		readToIndex -= m
	}
	return r, nil
}

func parseRequestLine(content string) (*RequestLine, error) {

	// if !strings.Contains(content, "HTTP/1.1") {
	// 	msg := "We only support HTTP/1.1 for now"
	// 	return nil, fmt.Errorf("%s", msg)
	// }

	c := strings.Split(content, " ")
	if len(c) != 3 {
		msg := "Request line is not correct"
		return nil, fmt.Errorf("%s", msg)
	}

	version := strings.Split(c[2], "/")[1]

	r := &RequestLine{
		HttpVersion:   version,
		RequestTarget: c[1],
		Method:        c[0],
	}

	return r, nil
}
