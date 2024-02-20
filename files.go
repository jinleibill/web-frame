package web_frame

import (
	lru "github.com/hashicorp/golang-lru"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FileUpload struct {
	FileField   string
	DstPathFunc func(*multipart.FileHeader) string
}

func (u *FileUpload) Handle() HandleFunc {
	return func(ctx *Context) {
		file, fileHeader, err := ctx.Req.FormFile(u.FileField)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败: " + err.Error())
			return
		}
		defer file.Close()
		dst := u.DstPathFunc(fileHeader)
		dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o666)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败: " + err.Error())
			return
		}
		defer dstFile.Close()
		_, err = io.CopyBuffer(dstFile, file, nil)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败: " + err.Error())
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("上传成功")
	}
}

type FileDownloader struct {
	Dir string
}

func (d *FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		req, err := ctx.QueryValue("file")
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("找不到目标文件参数")
			return
		}
		dst := filepath.Join(d.Dir, req)
		fn := filepath.Base(dst)
		dst, _ = filepath.Abs(dst)
		if !strings.Contains(dst, d.Dir) {
			ctx.RespStatusCode = http.StatusForbidden
			ctx.RespData = []byte("无权访问")
			return
		}
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")

		http.ServeFile(ctx.Resp, ctx.Req, dst)
	}
}

type StaticResourceHandlerOption func(handler *StaticResourceHandler)

type StaticResourceHandler struct {
	dir               string
	cache             *lru.Cache
	extContentTypeMap map[string]string
	maxSize           int
}

func NewStaticResourceHandler(dir string, opts ...StaticResourceHandlerOption) (*StaticResourceHandler, error) {
	c, err := lru.New(1000)
	if err != nil {
		return nil, err
	}
	res := &StaticResourceHandler{
		dir:   dir,
		cache: c,
		extContentTypeMap: map[string]string{
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "image/pdf",
		},
		maxSize: 1024 * 1024 * 10,
	}

	for _, opt := range opts {
		opt(res)
	}

	return res, nil
}

func (s *StaticResourceHandler) Handle(ctx *Context) {
	file, err := ctx.PathValue("file")
	if err != nil {
		ctx.RespStatusCode = http.StatusBadRequest
		ctx.RespData = []byte("请求路径错误")
		return
	}
	dst := filepath.Join(s.dir, file)
	dst, _ = filepath.Abs(dst)
	if !strings.Contains(dst, s.dir) {
		ctx.RespStatusCode = http.StatusForbidden
		ctx.RespData = []byte("无权访问")
		return
	}
	ext := filepath.Ext(dst)[1:]
	header := ctx.Resp.Header()
	if data, ok := s.cache.Get(file); ok {
		header.Set("Content-Type", s.extContentTypeMap[ext])
		header.Set("Content-Length", strconv.Itoa(len(data.([]byte))))

		ctx.RespData = data.([]byte)
		ctx.RespStatusCode = http.StatusOK
		return
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("服务器内部错误")
		return
	}
	if len(data) <= s.maxSize {
		s.cache.Add(file, data)
	}
	header.Set("Content-Type", s.extContentTypeMap[ext])
	header.Set("Content-Length", strconv.Itoa(len(data)))

	ctx.RespData = data
	ctx.RespStatusCode = http.StatusOK
}

func StaticWithMaxFileSize(maxSize int) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.maxSize = maxSize
	}
}

func StaticWithCache(c *lru.Cache) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.cache = c
	}
}

func StaticWithMoreExtension(extMap map[string]string) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		for ext, contextType := range extMap {
			handler.extContentTypeMap[ext] = contextType
		}
	}
}
