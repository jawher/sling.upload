package upload_test

import (
	"testing"

	"strings"

	"io/ioutil"

	"bytes"

	"os"

	"github.com/dghubble/sling"
	"github.com/jawher/sling.upload"
	"github.com/stretchr/testify/require"
)

func TestUpload(t *testing.T) {
	b := upload.New(
		upload.File("file", "fixtures/file.txt"),
		upload.Reader("reader", "reader.txt", strings.NewReader("readervalue")),
		upload.Param("param", "value"),
	)

	upload.New(
		upload.Part("paramName", upload.Filev("fixtures/file.txt"), upload.PartConfig{
			ContentType: "application/json",
		}),
		upload.File("paramName", upload.Reader(strings.NewReader("lol sdfasdf")), upload.PartConfig{
			Filename:    "data.json",
			ContentType: "application/json",
		}),
	)

	require.Equal(t, "multipart/form-data; boundary=SlingFormBoundary0amF3aGVy", b.ContentType())

	reader, err := b.Body()
	require.NoError(t, err)

	actualContent, err := ioutil.ReadAll(reader)
	require.NoError(t, err)

	ioutil.WriteFile("fixtures/actual.txt", actualContent, os.ModePerm)

	expectedContent, err := ioutil.ReadFile("fixtures/expected.txt")
	expectedContent = bytes.Replace(expectedContent, []byte("\n"), []byte("\r\n"), -1)
	require.NoError(t, err)

	require.Equal(t, expectedContent, actualContent)
}

func TestCompatibleWithSling(t *testing.T) {
	sling.New().BodyProvider(upload.New())
}
