package upload

import (
	"io"

	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"

	"fmt"

	"github.com/dghubble/sling"
)

var (
	multipartFromBoundary = "SlingFormBoundary0amF3aGVy"
)

// Part represents a part in a multipart request
type Part struct {
	// The part (or field) name
	Name string
	// The filename (optional)
	FileName string
	// The part content type
	ContentType string
	// The part content
	Content PartContent
	// Additional headers to add to the part
	Headers textproto.MIMEHeader
}

// PartContent provides the part's value to be serialized in the outgoing request
type PartContent interface {
	// Return the part's value or an error
	Get() (io.ReadCloser, error)
}

// String uses the provided string value as the part's content
func String(value string) PartContent {
	return stringContent(value)
}

type stringContent string

func (s stringContent) Get() (io.ReadCloser, error) {
	return ioutil.NopCloser(strings.NewReader(string(s))), nil
}

// File uses the provided file content as the part's value
func File(path string) PartContent {
	return fileContent(path)
}

type fileContent string

func (path fileContent) Get() (io.ReadCloser, error) {
	f, err := os.Open(string(path))
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Reader uses the provided reader as the part's value
func Reader(r io.Reader) PartContent {
	rc, ok := r.(io.ReadCloser)
	if !ok {
		rc = ioutil.NopCloser(r)
	}
	return &readerContent{
		r: rc,
	}
}

type readerContent struct {
	r io.ReadCloser
}

func (r *readerContent) Get() (io.ReadCloser, error) {
	return r.r, nil
}

type multipartBodyProvider []Part

// New creates a Sling body provider which serializes the provided parts to a form/multipart request
func New(parts ...Part) sling.BodyProvider {
	return multipartBodyProvider(parts)
}

func (p multipartBodyProvider) ContentType() string {
	return "multipart/form-data; boundary=" + multipartFromBoundary
}

func (p multipartBodyProvider) Body() (io.Reader, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	w.SetBoundary(multipartFromBoundary)

	for _, part := range p {
		h := part.Headers
		if h == nil {
			h = textproto.MIMEHeader{}
		}
		if part.ContentType != "" {
			h.Set("Content-Type", part.ContentType)
		}

		contentDisposition := fmt.Sprintf(`form-data; name="%s"`, escapeQuotes(part.Name))

		if part.FileName != "" {
			contentDisposition += fmt.Sprintf(`; filename="%s"`, escapeQuotes(part.FileName))
		}

		h.Set("Content-Disposition", contentDisposition)

		pw, err := w.CreatePart(h)
		if err != nil {
			return nil, err
		}

		if part.Content == nil {
			continue
		}
		reader, err := part.Content.Get()
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(pw, reader)
		reader.Close()
		if err != nil {
			return nil, err
		}

	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return &body, nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
