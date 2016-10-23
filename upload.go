package upload

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"

	"net/textproto"

	"github.com/dghubble/sling"
)

var (
	multipartFromBoundary = "SlingFormBoundary0amF3aGVy"
)

type Configurer func(*PartConfig)

func Header(key, value string) Configurer {
	return func(p *PartConfig) {
		p.Headers.Add(key, value)
	}
}

type PartConfig struct {
	Name     string
	Filename string
	Content  ContentProvider
	Headers  textproto.MIMEHeader
}

type PartHandler func(w *multipart.Writer) error

func Part(name string, content ContentProvider, configurers ...Configurer) PartHandler {
	pc := &PartConfig{
		Name:     name,
		Filename: "",
		Content:  content,
		Headers:  make(textproto.MIMEHeader),
	}

	for _, cfg := range configurers {
		cfg(pc)
	}

}

func defaultPartHandler(pc *PartConfig) PartHandler {
	return func(w *multipart.Writer) error {
		disposition := pc.Headers.Get("Content-Disposition")

		if disposition == "" {

		}
		w.CreatePart()
	}
}

// File creates a part from an existing file on disk. The name argument sets the part name, and path must point to an existing file on disk.
func File(name string, path string) Part {
	return func(w *multipart.Writer) error {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		fw, err := w.CreateFormFile(name, filepath.Base(path))
		if err != nil {
			return err
		}
		if _, err = io.Copy(fw, f); err != nil {
			return err
		}
		return err
	}
}

// Reader creates a part from a `io.Reader`. The name argument sets the part name, filename the uploaded file name and reader the content to be uploaded.
func Reader(name string, fileName string, reader io.Reader) Part {
	rc, ok := reader.(io.ReadCloser)
	if !ok {
		rc = ioutil.NopCloser(reader)
	}

	return func(w *multipart.Writer) error {
		defer rc.Close()
		fw, err := w.CreateFormFile(name, fileName)
		if err != nil {
			return err
		}
		if _, err = io.Copy(fw, reader); err != nil {
			return err
		}
		return nil
	}
}

// Param creates a non-file part (regular key-value). The name argument sets the part name, and path must point to an existing file on disk.
func Param(name string, value string) Part {
	return func(w *multipart.Writer) error {
		fw, err := w.CreateFormField(name)
		if err != nil {
			return err
		}
		if _, err := fw.Write([]byte(value)); err != nil {
			return err
		}
		return nil
	}
}

type multipartBodyProvider []Part

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
		if err := part(w); err != nil {
			return nil, err
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return &body, nil
}
