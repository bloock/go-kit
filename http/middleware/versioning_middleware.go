package middleware

import (
	"bytes"
	"encoding/json"
	pinned "github.com/bloock/go-kit/http/versioning"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func HandlerVersioning(vm *pinned.VersionManager, versions []*pinned.Version) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		originalWriter := ctx.Writer
		customResponseWriter := &customResponseWriter{ResponseWriter: ctx.Writer, bodyBuffer: bytes.NewBuffer(nil)}

		ctx.Writer = customResponseWriter
		version := getVersion(ctx, vm, versions)
		if ctx.Request.Body == nil {
			ctx.Next()
			return
		}

		if ctx.Request.Method != http.MethodGet {
			baseRequest := map[string]interface{}{}
			err := ctx.BindJSON(&baseRequest)
			if err != nil {
				return
			}

			applyRequest, err := vm.ApplyRequest(baseRequest, version, versions)
			if err != nil {
				_ = ctx.Error(err)
				ctx.Abort()
			}
			if applyRequest != nil {

				requestBytes, err := json.Marshal(applyRequest)
				if err != nil {
					_ = ctx.Error(err)
					return
				}
				ctx.Request.Body = io.NopCloser(bytes.NewReader(requestBytes))
			}
		}

		ctx.Next()

		status := ctx.Writer.Status()
		if len(ctx.Errors) > 0 {
			ctx.Writer = originalWriter
			return
		}

		b := customResponseWriter.bodyBuffer.Bytes()

		baseResponse := map[string]interface{}{}

		err := json.Unmarshal(b, &baseResponse)
		if err != nil {
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}
		response, err := vm.ApplyResponse(baseResponse, version, versions)
		if err != nil {
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}
		if response == nil {
			_, err = customResponseWriter.ResponseWriter.Write(b)
			return
		}
		respBytes, err := json.Marshal(response)
		if err != nil {
			_ = ctx.Error(err)
			return
		}
		ctx.Status(status)

		_, err = customResponseWriter.ResponseWriter.Write(respBytes)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

	}
}

func getVersion(ctx *gin.Context, vm *pinned.VersionManager, versions []*pinned.Version) *pinned.Version {
	version, err := vm.Parse(ctx.Request, versions)
	if err == pinned.ErrNoVersionSupplied || err == pinned.ErrInvalidVersion {
		version = vm.Oldest(versions)
	} else if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return nil
	}
	ctx.Set("version", version)
	return version
}

type customResponseWriter struct {
	gin.ResponseWriter
	bodyBuffer *bytes.Buffer
	responses  []byte
}

func (w *customResponseWriter) Write(data []byte) (int, error) {
	return w.bodyBuffer.Write(data)
}

func (w *customResponseWriter) WriteString(s string) (int, error) {
	return w.bodyBuffer.WriteString(s)
}

func (w *customResponseWriter) Status() int {
	return w.ResponseWriter.Status()
}

func (w *customResponseWriter) Body() string {
	return w.bodyBuffer.String()
}

type VersionedHandler interface {
	Handler() gin.HandlerFunc
	Versions() []*pinned.Version
}
