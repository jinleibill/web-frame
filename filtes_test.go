package web_frame

import (
	"github.com/stretchr/testify/require"
	"html/template"
	"log"
	"mime/multipart"
	"path/filepath"
	"testing"
)

func TestUpload(t *testing.T) {
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	require.NoError(t, err)
	engine := &GoTemplateEngine{
		T: tpl,
	}
	h := NewHTTPServer(ServerWithTemplateEngine(engine))
	h.Get("/upload", func(ctx *Context) {
		err := ctx.Render("upload.gohtml", nil)
		if err != nil {
			log.Println(err)
		}
	})
	fu := FileUpload{
		FileField: "myfile",
		DstPathFunc: func(header *multipart.FileHeader) string {
			return filepath.Join("testdata", "upload", header.Filename)
		},
	}
	h.Post("/upload", fu.Handle())

	_ = h.Start(":8081")
}

func TestDownload(t *testing.T) {
	h := NewHTTPServer()
	fu := &FileDownloader{
		Dir: filepath.Join("testdata", "download"),
	}
	h.Get("/download", fu.Handle())

	_ = h.Start(":8081")
}

func TestStaticResourceHandler_Handle(t *testing.T) {
	h := NewHTTPServer()
	fu, err := NewStaticResourceHandler(filepath.Join("testdata", "static"))
	require.NoError(t, err)
	h.Get("/static/:file", fu.Handle)

	_ = h.Start(":8081")
}
