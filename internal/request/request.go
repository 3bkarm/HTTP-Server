package request

import (
	"bytes"
	"errors"
	"httpserver/internal/headers"
	"io"
)

var SEPARATOR = []byte("\r\n")

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(b []byte) (rl *RequestLine, read int, er error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	requestLine := b[:idx]
	read = idx + len(SEPARATOR)

	parts := bytes.Split(requestLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, errors.New("malformed request-line")
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, errors.New("unsupported http version")
	}

	rl = &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}
	return rl, read, nil
}

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
	StateDone    parserState = "done"
	StateError   parserState = "error"
)

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	Body        string
	State       parserState
}

// TODO: make it static like
func newRequest() *Request {
	return &Request{
		State:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

func (rq *Request) done() bool {
	return rq.State == StateDone
}

func (rq *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentData := data[read:]
		switch rq.State {
		case StateError:
			return 0, errors.New("request in error state")
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				rq.State = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			rq.RequestLine = *rl
			rq.State = StateHeaders
			read += n
		case StateHeaders:
			n, done, err := rq.Headers.Parse(currentData)
			if err != nil {
				rq.State = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			read += n
			if done {
				rq.State = StateBody
			}
		case StateBody:
			bodyLength := rq.Headers.GetInt("content-length", 0)
			if bodyLength == 0 {
				rq.State = StateDone
				break outer
			}
			remainingBody := min(bodyLength-len(rq.Body), len(currentData))
			if remainingBody == 0 {
				break outer
			}
			rq.Body += string(currentData[:remainingBody])
			read += remainingBody
			if len(rq.Body) == bodyLength {
				rq.State = StateDone
			}
		case StateDone:
			break outer
		default:
			panic("missing request state case!")
		}
	}
	return read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	// NOTE: buffer size may be not enough
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		// TODO: better explaining for the error
		if err != nil {
			return nil, err
		}
		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}
	return request, nil
}
