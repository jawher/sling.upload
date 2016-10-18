package upload

import (
	"testing"

	"strings"

	"io/ioutil"

	"bytes"

	"os"

	"github.com/dghubble/sling"
	"github.com/stretchr/testify/require"
)

func TestUpload(t *testing.T) {
	b := New(
		File("file", "fixtures/file.txt"),
		Reader("reader", "reader.txt", strings.NewReader("readervalue")),
		Param("param", "value"),
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
	sling.New().BodyProvider(New())
}
