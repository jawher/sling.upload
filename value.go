package upload

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ContentProvider func() (io.Reader, string, error)

func Filev(path string) Value {
	return func() (io.Reader, string, error) {
		f, err := os.Open(path)
		if err != nil {
			return nil, "", err
		}
		return f, filepath.Base(path), nil
	}
}

func Stringv(s string) Value {
	return func() (io.Reader, string, error) {
		return strings.NewReader(s), "", nil
	}
}

func Bytesv(data []byte) Value {
	return func() (io.Reader, string, error) {
		return bytes.NewReader(data), "", nil
	}
}
