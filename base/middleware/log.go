package middleware

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	logRequestBodyLimit = 10240
	logResponseLimit    = 10240
)

type responseWriterWrapper struct {
	original gin.ResponseWriter
	resp     strings.Builder
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	w.original.WriteHeader(code)
}

func (w *responseWriterWrapper) WriteHeaderNow() {
	w.original.WriteHeaderNow()
}

func (w *responseWriterWrapper) Write(data []byte) (n int, err error) {
	w.resp.Write(data)
	return w.original.Write(data)
}

func (w *responseWriterWrapper) WriteString(s string) (n int, err error) {
	w.resp.WriteString(s)
	return w.original.WriteString(s)
}

func (w *responseWriterWrapper) Status() int {
	return w.original.Status()
}

func (w *responseWriterWrapper) Size() int {
	return w.original.Size()
}

func (w *responseWriterWrapper) Written() bool {
	return w.original.Written()
}

func (w *responseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.original.Hijack()
}

func (w *responseWriterWrapper) CloseNotify() <-chan bool {
	return w.original.CloseNotify()
}

func (w *responseWriterWrapper) Flush() {
	w.original.Flush()
}

func (w *responseWriterWrapper) Pusher() (pusher http.Pusher) {
	return w.original.Pusher()
}

func (w *responseWriterWrapper) Header() http.Header {
	return w.original.Header()
}

func cloneRequestBodyForLog(c *gin.Context) string {
	contentType := c.GetHeader("Content-Type")
	if contentType != "" && !strings.Contains(contentType, "application/json") {
		return contentType
	}
	length := c.GetHeader("Content-Length")
	val64, err := strconv.ParseInt(length, 10, 0)
	if err != nil {
		return fmt.Sprintf("parse Content-Length err: %s", length)
	}
	if val64 > logRequestBodyLimit {
		return fmt.Sprintf("huge body [%d]", val64)
	}
	body, _ := c.GetRawData()
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return string(body)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startAt := time.Now()
		body := cloneRequestBodyForLog(c)
		writerWrapper := &responseWriterWrapper{original: c.Writer}
		c.Writer = writerWrapper
		c.Next()
		cost := time.Now().Sub(startAt)
		resp := writerWrapper.resp.String()
		if len(resp) > logResponseLimit {
			resp = resp[:logResponseLimit] + "..."
		}
		zap.L().Info(
			fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.String()),
			zap.String("ip", c.ClientIP()),
			zap.String("time", cost.String()),
			zap.String("request", body),
			zap.Any("header", c.Request.Header),
			zap.String("response", resp),
		)
	}
}
