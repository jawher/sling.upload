package upload_test

import (
	"testing"

	"strings"

	"io/ioutil"

	"os"

	"flag"

	"github.com/dghubble/sling"
	upload "github.com/jawher/sling.upload"
	"github.com/stretchr/testify/require"
)

var writeActual bool

func init() {
	flag.BoolVar(&writeActual, "w", false, "Write the actual generated body to fixtures/actual.txt")
	flag.Parse()
}

func TestUpload(t *testing.T) {

	b := upload.New(
		upload.Part{
			Name:        "file",
			FileName:    "file.txt",
			Content:     upload.File("fixtures/file.txt"),
			ContentType: "application/json",
		},
		upload.Part{
			Name:    "reader",
			Content: upload.Reader(strings.NewReader("readervalue")),
		},
		upload.Part{
			Name:    "param",
			Content: upload.String("value"),
		},
	)

	require.Equal(t, "multipart/form-data; boundary=SlingFormBoundary0amF3aGVy", b.ContentType())

	reader, err := b.Body()
	require.NoError(t, err)

	actualContent, err := ioutil.ReadAll(reader)
	require.NoError(t, err)

	if writeActual {
		ioutil.WriteFile("fixtures/actual.txt", actualContent, os.ModePerm)
	}

	expectedContent, err := ioutil.ReadFile("fixtures/expected.txt")
	require.NoError(t, err)

	require.Equal(t, expectedContent, actualContent)
}

func TestCompatibleWithSling(t *testing.T) {
	sling.New().BodyProvider(upload.New())
}

func Example() {
	sling.New().Post("http://localhost:4000/upload").BodyProvider(
		upload.New(
			upload.Part{
				Name:        "file",
				FileName:    "file.json",
				Content:     upload.File("~/file.json"),
				ContentType: "application/json",
			},
			upload.Part{
				Name:    "name",
				Content: upload.String("sling.upload"),
			},
		),
	).Receive(nil, nil)
}
