package xin

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const defaultMultipartMemory = 32 << 20 // 32 MB

var MaxMultipartMemory int64 = defaultMultipartMemory

// get is an internal method and returns a map which satisfy conditions.
func get(m map[string][]string, key string) (map[string]string, bool) {
	dicts := make(map[string]string)
	exist := false
	for k, v := range m {
		if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
			if j := strings.IndexByte(k[i+1:], ']'); j >= 1 {
				exist = true
				dicts[k[i+1:][:j]] = v[0]
			}
		}
	}
	return dicts, exist
}

// Query returns the keyed url query value if it exists,
// otherwise it returns an empty string `("")`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
//     GET /path?id=1234&name=Manu&value=
// 	   c.Query("id") == "1234"
// 	   c.Query("name") == "Manu"
// 	   c.Query("value") == ""
// 	   c.Query("wtf") == ""
func Query(r *http.Request, key string) string {
	value, _ := GetQuery(r, key)
	return value
}

// DefaultQuery returns the keyed url query value if it exists,
// otherwise it returns the specified defaultValue string.
// See: Query() and GetQuery() for further information.
//     GET /?name=Manu&lastname=
//     c.DefaultQuery("name", "unknown") == "Manu"
//     c.DefaultQuery("id", "none") == "none"
//     c.DefaultQuery("lastname", "none") == ""
func DefaultQuery(r *http.Request, key, defaultValue string) string {
	if value, ok := GetQuery(r, key); ok {
		return value
	}
	return defaultValue
}

// GetQuery is like Query(), it returns the keyed url query value
// if it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns `("", false)`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
//     GET /?name=Manu&lastname=
//     ("Manu", true) == c.GetQuery("name")
//     ("", false) == c.GetQuery("id")
//     ("", true) == c.GetQuery("lastname")
func GetQuery(r *http.Request, key string) (string, bool) {
	if values, ok := GetQueryArray(r, key); ok {
		return values[0], ok
	}
	return "", false
}

// QueryArray returns a slice of strings for a given query key.
// The length of the slice depends on the number of params with the given key.

func QueryArray(r *http.Request, key string) []string {
	values, _ := GetQueryArray(r, key)
	return values
}

// GetQueryArray returns a slice of strings for a given query key, plus
// a boolean value whether at least one value exists for the given key.
func GetQueryArray(r *http.Request, key string) ([]string, bool) {
	if values, ok := r.URL.Query()[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

// QueryMap returns a map for a given query key.
func QueryMap(r *http.Request, key string) map[string]string {
	dicts, _ := GetQueryMap(r, key)
	return dicts
}

// GetQueryMap returns a map for a given query key, plus a boolean value
// whether at least one value exists for the given key.
func GetQueryMap(r *http.Request, key string) (map[string]string, bool) {
	return get(r.URL.Query(), key)
}

// PostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns an empty string `("")`.
func PostForm(r *http.Request, key string) string {
	value, _ := GetPostForm(r, key)
	return value
}

// DefaultPostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns the specified defaultValue string.
// See: PostForm() and GetPostForm() for further information.
func DefaultPostForm(r *http.Request, key, defaultValue string) string {
	if value, ok := GetPostForm(r, key); ok {
		return value
	}
	return defaultValue
}

// GetPostForm is like PostForm(key). It returns the specified key from a POST urlencoded
// form or multipart form when it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns ("", false).
// For example, during a PATCH request to update the user's email:
//     email=mail@example.com  -->  ("mail@example.com", true) := GetPostForm("email") // set email to "mail@example.com"
// 	   email=                  -->  ("", true) := GetPostForm("email") // set email to ""
//                             -->  ("", false) := GetPostForm("email") // do nothing with email
func GetPostForm(r *http.Request, key string) (string, bool) {
	if values, ok := GetPostFormArray(r, key); ok {
		return values[0], ok
	}
	return "", false
}

// PostFormArray returns a slice of strings for a given form key.
// The length of the slice depends on the number of params with the given key.
func PostFormArray(r *http.Request, key string) []string {
	values, _ := GetPostFormArray(r, key)
	return values
}

// GetPostFormArray returns a slice of strings for a given form key, plus
// a boolean value whether at least one value exists for the given key.
func GetPostFormArray(r *http.Request, key string) ([]string, bool) {
	if err := r.ParseMultipartForm(MaxMultipartMemory); err != nil {
		//if err != http.ErrNotMultipart {
		//	debugPrint("error on parse multipart form array: %v", err)
		//}
	}
	if values := r.PostForm[key]; len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

// PostFormMap returns a map for a given form key.
func PostFormMap(r *http.Request, key string) map[string]string {
	dicts, _ := GetPostFormMap(r, key)
	return dicts
}

// GetPostFormMap returns a map for a given form key, plus a boolean value
// whether at least one value exists for the given key.
func GetPostFormMap(r *http.Request, key string) (map[string]string, bool) {
	if err := r.ParseMultipartForm(MaxMultipartMemory); err != nil {
		//if err != http.ErrNotMultipart {
		//	debugPrint("error on parse multipart form map: %v", err)
		//}
	}
	return get(r.PostForm, key)
}

// FormFile returns the first file for the provided form key.
func FormFile(r *http.Request, name string) (*multipart.FileHeader, error) {
	if r.MultipartForm == nil {
		if err := r.ParseMultipartForm(MaxMultipartMemory); err != nil {
			return nil, err
		}
	}
	_, fh, err := r.FormFile(name)
	return fh, err
}

// MultipartForm is the parsed multipart form, including file uploads.
func MultipartForm(r *http.Request) (*multipart.Form, error) {
	err := r.ParseMultipartForm(MaxMultipartMemory)
	return r.MultipartForm, err
}

// SaveUploadedFile uploads the form file to specific dst.
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
