package errorx

import "strconv"

const defaultCode = 1001

type CodeError struct {
	Code        int           `json:"code"`
	Msg         string        `json:"msg"`
	LangCode    string        `json:"langCode"`
	Placeholder []interface{} `json:"placeholder"`
}

type CodeErrorResponse struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	LangCode string `json:"langCode"`
}

func NewRawCodeError(code int, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg, LangCode: strconv.Itoa(code), Placeholder: make([]interface{}, 0)}
}

func NewRawCodeErrorWithLangCode(code int, msg, langCode string) *CodeError {
	return &CodeError{Code: code, Msg: msg, LangCode: langCode, Placeholder: make([]interface{}, 0)}
}

func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg, Placeholder: make([]interface{}, 0)}
}

func NewCodeErrorWithLang(code int, msg string, langCode string) error {
	return &CodeError{Code: code, Msg: msg, LangCode: langCode, Placeholder: make([]interface{}, 0)}
}

func NewErrorFromCode(c *CodeError) error {
	return &CodeError{Code: c.Code, Msg: c.Msg, LangCode: c.LangCode, Placeholder: make([]interface{}, 0)}
}

func NewErrorf(c *CodeError, v ...interface{}) error {
	placeholder := make([]interface{}, 0)
	for _, pv := range v {
		placeholder = append(placeholder, pv)
	}
	return &CodeError{Code: c.Code, Msg: c.Msg, LangCode: c.LangCode, Placeholder: placeholder}
}

func NewDefaultError(msg string) error {
	return NewCodeError(defaultCode, msg)
}

func (e *CodeError) Error() string {
	return e.Msg
}

func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code:     e.Code,
		Msg:      e.Msg,
		LangCode: e.LangCode,
	}
}
