package request

import (
	"HttpToTcp/internal/headers"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const bufferSize = 8

type state int

const (
	initialized state = iota
	done
	requestStateParsingHeaders
	requestStateParsingBody
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
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
	if !strings.Contains(string(data), "\r\n") { // Very Important:
		return 0, nil
	}

	//Depending on the numberBytesPerRead we can recieve data
	//that exceeds \r\n so we only treat data before that (split)
	//the lefovers will be shifted in the buffer
	if r.state == initialized {
		n, startLine, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		r.RequestLine = *startLine
		r.state = requestStateParsingHeaders

		return n, nil
	} else if r.state == requestStateParsingHeaders {
		n, ok, err := r.Headers.Parse([]byte(string(data)))
		if err != nil {
			return 0, err
		}
		if ok {
			r.state = requestStateParsingBody
		}
		return n, nil
	} else {
		content := r.Headers.Get("Content-Length")
		if content != "" {
			contentLenght, err := strconv.Atoi(content)
			if err != nil {
				r.state = done
				return 0, err
			}

			fmt.Println(">", contentLenght, len(data[2:]))
			if contentLenght == len(data[2:]) {
				r.Body = data[2:]
				r.state = done
				return len(data), nil
			} else {
				r.state = done
				msg := "Invalid Content-Length"
				return 0, fmt.Errorf("%s", msg)
			}
		}
		r.state = done
		return 0, nil
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	r := &Request{
		state: initialized,
	}
	r.Headers = headers.NewHeaders()

	for r.state != done {
		if len(buff) == cap(buff) {
			buff2 := make([]byte, len(buff)*2, cap(buff)*2)
			copy(buff2, buff)
			buff = buff2
		}

		// Read cr.numBytesPerRead into our buffer
		// Return number of bytes copied!! to the buffer
		numBytesRead, err := reader.Read(buff[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.state = done
			}
			return nil, err
		}

		readToIndex += numBytesRead

		numBytesParsed, err := r.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}

		// Shift leftover data to front of buffer
		if numBytesParsed > 0 {
			copy(buff, buff[numBytesParsed:readToIndex])
			readToIndex -= numBytesParsed
		}
	}
	return r, nil
}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	c := string(data)

	if !strings.Contains(c, "HTTP/1.1") {
		msg := "We only support HTTP/1.1 for now"
		return 0, nil, fmt.Errorf("%s", msg)
	}

	c = strings.Split(c, "\r\n")[0] // [1] will be leftovers
	l := strings.Fields(c)
	if len(l) != 3 {
		msg := "Request line is not correct"
		return 0, nil, fmt.Errorf("%s", msg)
	}

	version := strings.Split(l[2], "/")[1]

	r := &RequestLine{
		HttpVersion:   version,
		RequestTarget: l[1],
		Method:        l[0],
	}

	return len(c) + len("\r\n"), r, nil
}
