# Sling.Upload

[![Build Status](https://travis-ci.org/jawher/sling.upload.svg?branch=master)](https://travis-ci.org/jawher/sling.upload)
[![GoDoc](https://godoc.org/github.com/jawher/sling.upload?status.svg)](https://godoc.org/github.com/jawher/sling.upload)

This packages provides a multipart [BodyProvider](https://godoc.org/github.com/dghubble/sling#BodyProvider) implementation to be used with the [Sling HTTP client](https://github.com/dghubble/sling).

## Installation


```
go get github.com/jawher/sling.upload
```

## Usage

To create an instance of the multipart body provider, call the `upload.New` function passing in the list of parts to upload.

For example:

```go
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
```

