package charset

import (
	"io"
	"golang.org/x/net/html/charset"
)

func NewReader(r io.Reader, contentType string) (io.Reader, error) {
	return charset.NewReader(r,contentType)
}
