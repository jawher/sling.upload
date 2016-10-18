# Sling.Upload

This packages provides a multipart [BodyProvider](https://godoc.org/github.com/dghubble/sling#BodyProvider) implementation to be used with the [Sling HTTP client](https://github.com/dghubble/sling).

## Installation

A simple go get:

```
go get github.com/jawher/sling.upload
```

## Usage

To create an instance of the multipart body provider, call the `upload.New` function passing in the list of parts/files to upload.

For example:

```go
sling.New().Post("http://localhost:4000/upload").BodyProvider(
	upload.New(
		upload.File("file1", "~/data.json"),
		upload.Reader("file2", "data.txt", strings.NewReader("Text payload")),
		upload.Param("to", "SHOP"),
	)).Receive(nil, nil)
```

## Supported parts

### Files

Use `upload.File` to construct a part from an existing file in disk.