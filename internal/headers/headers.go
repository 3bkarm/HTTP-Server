package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var registerNurse = []byte("\r\n")

type Headers struct {
	headers map[string]string
}

// TODO: make it like static
func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func GetDefaultHeaders(contentLength int) *Headers {
	h := NewHeaders()
	h.Set("content-length", fmt.Sprintf("%d", contentLength))
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")
	return h
}

func isToken(str string) bool {
	for _, ch := range str {
		ok := false
		if (ch >= 'A' && ch <= 'Z') ||
			(ch >= 'a' && ch <= 'z') ||
			(ch >= '0' && ch <= '9') {
			ok = true
		}
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			ok = true
		}
		if !ok {
			return false
		}
	}
	return true
}

func (h *Headers) Get(key string) (string, bool) {
	str, ok := h.headers[strings.ToLower(key)]
	return str, ok
}

func (h *Headers) Set(key string, value string) {
	key = strings.ToLower(key)
	if v, ok := h.headers[key]; ok {
		h.headers[key] = fmt.Sprintf("%s, %s", v, value)
	} else {
		h.headers[key] = value
	}
}

func (h *Headers) Replace(key string, value string) {
	h.headers[key] = value
}

func (h *Headers) ForEach(apply func(k, v string)) {
	for k, v := range h.headers {
		apply(k, v)
	}
}

func (h *Headers) GetInt(k string, defaultValue int) int {
	valueStr, exists := h.Get(k)
	if !exists {
		return defaultValue
	}
	// TODO: maybe better dealing with invalid data value
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", errors.New("malformed field line")
	}

	value := bytes.TrimSpace(parts[1])
	if bytes.HasSuffix(parts[0], []byte(" ")) {
		return "", "", errors.New("malformed field name")
	}
	name := bytes.TrimSpace(parts[0])

	return string(name), string(value), nil
}

func (h *Headers) Parse(data []byte) (read int, done bool, err error) {
	read = 0
	done = false
	for {
		idx := bytes.Index(data[read:], registerNurse)
		if idx == -1 {
			break
		}

		if idx == 0 {
			done = true
			read += len(registerNurse)
			break
		}

		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}

		if !isToken(name) {
			return 0, false, errors.New("malformed header name")
		}

		h.Set(name, value)
		read += idx + len(registerNurse)
	}
	return read, done, nil
}
