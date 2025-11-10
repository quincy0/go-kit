package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"go-kit/core/mapping"
	"go-kit/rest/internal/encoding"
	"go-kit/rest/internal/header"
	"go-kit/rest/pathvar"
)

const (
	formKey           = "form"
	pathKey           = "path"
	maxMemory         = 32 << 20 // 32MB
	maxBodyLen        = 8 << 20  // 8MB
	separator         = ";"
	tokensInAttribute = 2
)

var (
	formUnmarshaler = mapping.NewUnmarshaler(formKey, mapping.WithStringValues())
	pathUnmarshaler = mapping.NewUnmarshaler(pathKey, mapping.WithStringValues())
)

var ErrUserInvalid = errors.New("user invalid")

// Parse parses the request.
func Parse(r *http.Request, v interface{}) error {
	if err := ParsePath(r, v); err != nil {
		return err
	}

	if err := ParseForm(r, v); err != nil {
		return err
	}

	if err := ParseHeaders(r, v); err != nil {
		return err
	}

	return ParseJsonBody(r, v)
}

// ParseHeaders parses the headers request.
func ParseHeaders(r *http.Request, v interface{}) error {
	return encoding.ParseHeaders(r.Header, v)
}

// ParseForm parses the form request.
func ParseForm(r *http.Request, v interface{}) error {
	params, err := GetFormValues(r)
	if err != nil {
		return err
	}

	return formUnmarshaler.Unmarshal(params, v)
}

// ParseHeader parses the request header and returns a map.
func ParseHeader(headerValue string) map[string]string {
	ret := make(map[string]string)
	fields := strings.Split(headerValue, separator)

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if len(field) == 0 {
			continue
		}

		kv := strings.SplitN(field, "=", tokensInAttribute)
		if len(kv) != tokensInAttribute {
			continue
		}

		ret[kv[0]] = kv[1]
	}

	return ret
}

// ParseJsonBody parses the post request which contains json in body.
func ParseJsonBody(r *http.Request, v interface{}) error {
	if withJsonBody(r) {
		reader := io.LimitReader(r.Body, maxBodyLen)
		return mapping.UnmarshalJsonReader(reader, v)
	}

	return mapping.UnmarshalJsonMap(nil, v)
}

// ParsePath parses the symbols reside in url path.
// Like http://localhost/bag/:name
func ParsePath(r *http.Request, v interface{}) error {
	vars := pathvar.Vars(r)
	m := make(map[string]interface{}, len(vars))
	for k, v := range vars {
		m[k] = v
	}

	return pathUnmarshaler.Unmarshal(m, v)
}

func withJsonBody(r *http.Request) bool {
	return r.ContentLength > 0 && strings.Contains(r.Header.Get(header.ContentType), header.ApplicationJson)
}

type RequestSvc interface {
	User() *User
	Header(key string) string
	UserJson() string
}

type requestSvc struct {
	user     *User
	header   map[string]string
	userJson string
}

func NewRequestWithTokenValue(r *http.Request, tokenValue string) (RequestSvc, error) {
	re := &requestSvc{
		header:   make(map[string]string),
		userJson: tokenValue,
	}

	for k, v := range r.Header {
		re.header[k] = v[0]
		re.header[strings.ToLower(k)] = v[0] // 转成小写，兼容老版本
	}

	if len(tokenValue) == 0 {
		return re, nil
	}

	acInfo := User{}
	if err := json.Unmarshal([]byte(tokenValue), &acInfo); err != nil {
		return re, err
	}

	re.user = &acInfo

	return re, nil
}

func (rs *requestSvc) User() *User {
	return rs.user
}

func (rs *requestSvc) Header(key string) string {
	val, ok := rs.header[key]
	if !ok {
		val, ok = rs.header[strings.ToLower(key)]
		if !ok {
			return ""
		}
	}
	return val
}

func (rs *requestSvc) UserJson() string {
	return rs.userJson
}
