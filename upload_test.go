package upload

import (
	"testing"

	"strings"

	"io/ioutil"

	"bytes"

	"github.com/dghubble/sling"
	"github.com/stretchr/testify/require"
)

func TestUpload(t *testing.T) {
	b := New(
		Param("param", "value"),
		Reader("reader", "reader.txt", strings.NewReader("readervalue")),
		File("file", "fixtures/file.txt"),
	)

	require.Equal(t, "multipart/form-data; boundary=SlingFormBoundary0amF3aGVy", b.ContentType())

	reader, err := b.Body()
	require.NoError(t, err)

	actualContent, err := ioutil.ReadAll(reader)
	require.NoError(t, err)

	expectedContent, err := ioutil.ReadFile("fixtures/expected.txt")
	expectedContent = bytes.Replace(expectedContent, []byte("\n"), []byte("\r\n"), -1)
	require.NoError(t, err)

	require.Equal(t, expectedContent, actualContent)
}

func TestCompatibleWithSling(t *testing.T) {
	sling.New().BodyProvider(New())
}
