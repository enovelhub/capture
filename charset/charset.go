package charset

import (
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
)

func NewReader(r io.Reader, contentType string) (io.Reader, error) {
	return charset.NewReader(r, contentType)
}

func ToUTF8(resp *http.Response) (io.Reader, error) {
	return NewReader(resp.Body, resp.Header.Get("Content-Type"))
}
