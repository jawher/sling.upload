package upload

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/dghubble/sling"
)

var (
	multipartFromBoundary = "SlingFormBoundary0amF3aGVy"
)

type multipartBodyProvider struct {
	files       map[string]string
	readers     map[string]namedReader
	extraParams map[string]string
}

type namedReader struct {
	fileName string
	reader   io.ReadCloser
}

type part func(*multipartBodyProvider)

// File creates a part from an existing file on disk. The name argument sets the part name, and path must point to an existing file on disk.
func File(name string, path string) part {
	return func(m *multipartBodyProvider) {
		m.files[name] = path
	}
}

// Reader creates a part from a `io.Reader`. The name argument sets the part name, filename the uploaded file name and reader the content to be uploaded.
func Reader(name string, fileName string, reader io.Reader) part {
	rc, ok := reader.(io.ReadCloser)
	if !ok {
		rc = ioutil.NopCloser(reader)
	}

	return func(m *multipartBodyProvider) {
		m.readers[name] = namedReader{
			fileName: fileName,
			reader:   rc,
		}
	}
}

// Param creates a non-file part (regular key-value). The name argument sets the part name, and path must point to an existing file on disk.
func Param(name string, value string) part {
	return func(m *multipartBodyProvider) {
		m.extraParams[name] = value
	}
}

func New(parts ...part) sling.BodyProvider {
	m := &multipartBodyProvider{
		files:       map[string]string{},
		readers:     map[string]namedReader{},
		extraParams: map[string]string{},
	}

	for _, c := range parts {
		c(m)
	}

	return m
}

func (p multipartBodyProvider) ContentType() string {
	return "multipart/form-data; boundary=" + multipartFromBoundary
}

func (p multipartBodyProvider) Body() (io.Reader, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	w.SetBoundary(multipartFromBoundary)

	for paramName, filePath := range p.files {
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		fw, err := w.CreateFormFile(paramName, filepath.Base(filePath))
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(fw, f); err != nil {
			return nil, err
		}
	}

	for paramName, namedReader := range p.readers {
		defer namedReader.reader.Close()
		fw, err := w.CreateFormFile(paramName, namedReader.fileName)
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(fw, namedReader.reader); err != nil {
			return nil, err
		}
	}
	// Add the other fields
	for paramName, paramValue := range p.extraParams {
		fw, err := w.CreateFormField(paramName)
		if err != nil {
			return nil, err
		}
		if _, err := fw.Write([]byte(paramValue)); err != nil {
			return nil, err
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return &body, nil
}
