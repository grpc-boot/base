package main

import (
	"bytes"
	_ "embed"
	"net/http"

	"github.com/grpc-boot/base/v2/utils"
)

//go:embed example.md
var example []byte

func main() {
	http.HandleFunc("/black", func(writer http.ResponseWriter, request *http.Request) {
		buf := bytes.NewBuffer(nil)

		buf.WriteString(`<!DOCTYPE html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.5.0/github-markdown-light.min.css" />
	</head>
	<body>
		<article class="markdown-body" style="padding:10px 30px;">`)
		buf.Write(utils.Markdown2Html(example))
		buf.WriteString(`</article>
	</body>
</html>`)
		_, _ = writer.Write(buf.Bytes())
	})

	http.ListenAndServe(":8080", http.DefaultServeMux)
}
