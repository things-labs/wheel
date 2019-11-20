// Package httutil is useful net/http
// almost copy from github.com/gin-gonic/gin for net/http used
package httputil

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin/render"
)

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

func Render(w http.ResponseWriter, code int, r render.Render) {
	w.WriteHeader(code)

	if !bodyAllowedForStatus(code) {
		r.WriteContentType(w)
		// TODO:
		return
	}

	if err := r.Render(w); err != nil {
		panic(err)
	}
}

func IndentedJSON(w http.ResponseWriter, code int, obj interface{}) {
	Render(w, code, render.IndentedJSON{Data: obj})
}

func SecureJSON(w http.ResponseWriter, code int, obj interface{}) {
	Render(w, code, render.SecureJSON{Prefix: "" /*TODO*/, Data: obj})
}

func JSONP(w http.ResponseWriter, code int, obj interface{}) {
	//callback := c.DefaultQuery("callback", "")
	//if callback == "" {
	//	Render(w, code, render.JSON{Data: obj})
	//	return
	//}
	//Render(w, code, render.JsonpJSON{Callback: callback, Data: obj})
}

func JSON(w http.ResponseWriter, code int, obj interface{}) {
	Render(w, code, render.JSON{Data: obj})
}

func AsciiJSON(w http.ResponseWriter, code int, obj interface{}) {
	Render(w, code, render.AsciiJSON{Data: obj})
}

func PureJSON(w http.ResponseWriter, code int, obj interface{}) {
	Render(w, code, render.PureJSON{Data: obj})
}

// YAML serializes the given struct as YAML into the response body.
func YAML(w http.ResponseWriter, code int, obj interface{}) {
	Render(w, code, render.YAML{Data: obj})
}

// ProtoBuf serializes the given struct as ProtoBuf into the response body.
func ProtoBuf(w http.ResponseWriter, code int, obj interface{}) {
	Render(w, code, render.ProtoBuf{Data: obj})
}

// String writes the given string into the response body.
func String(w http.ResponseWriter, code int, format string, values ...interface{}) {
	Render(w, code, render.String{Format: format, Data: values})
}

// Redirect returns a HTTP redirect to the specific location.
func Redirect(w http.ResponseWriter, r *http.Request, code int, location string) {
	Render(w, -1, render.Redirect{
		Code:     code,
		Location: location,
		Request:  r,
	})
}

// Data writes some data into the body stream and updates the HTTP code.
func Data(w http.ResponseWriter, code int, contentType string, data []byte) {
	Render(w, code, render.Data{
		ContentType: contentType,
		Data:        data,
	})
}

// DataFromReader writes the specified reader into the body stream and updates the HTTP code.
func DataFromReader(w http.ResponseWriter, code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string) {
	Render(w, code, render.Reader{
		Headers:       extraHeaders,
		ContentType:   contentType,
		ContentLength: contentLength,
		Reader:        reader,
	})
}
