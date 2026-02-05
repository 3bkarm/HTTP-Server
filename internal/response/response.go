package response

import (
	"errors"
	"fmt"
	"httpserver/internal/headers"
	"io"
)

type Response struct {
}

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func OkRespond() []byte {
	return []byte(`
	<html>
		<head>
			<title>200 OK</title>
		</head>
		<body>
			<h1>Success!</h1>
			<p>Your request was an absolute banger.</p>
		</body>
	</html>
	`)
}

func BadRequestRespond() []byte {
	return []byte(`
	<html>
		<head>
			<title>400 Bad Request</title>
		</head>
		<body>
			<h1>Bad Request</h1>
			<p>Your request honestly kinda sucked.</p>
		</body>
	</html>
	`)
}

func InternalServerErrorRespond() []byte {
	return []byte(`
	<html>
		<head>
			<title>500 Internal Server Error</title>
		</head>
		<body>
			<h1>Internal Server Error</h1>
			<p>Okay, you know what? This one is on me.</p>
		</body>
	</html>
	`)
}

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := []byte{}
	switch statusCode {
	case StatusOk:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusInternalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return errors.New("unrecognized status code")
	}
	_, err := w.writer.Write(statusLine)
	return err
}

func (w *Writer) WriteHeaders(h *headers.Headers) error {
	headerLines := []byte{}
	h.ForEach(func(k, v string) {
		headerLines = fmt.Appendf(headerLines, "%s: %s\r\n", k, v)
	})
	headerLines = fmt.Append(headerLines, "\r\n")
	_, err := w.writer.Write(headerLines)
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	return n, err
}
