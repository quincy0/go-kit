package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/quincy0/go-kit/core/logx"
	"github.com/quincy0/go-kit/rest/errorx"
	"github.com/quincy0/go-kit/rest/internal/errcode"
	"github.com/quincy0/go-kit/rest/internal/header"
	"github.com/quincy0/go-kit/rest/lang"
)

var (
	errorHandler func(error) (int, interface{})
	lock         sync.RWMutex
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Error writes err into w.
func Error(w http.ResponseWriter, err error, fns ...func(w http.ResponseWriter, err error)) {
	lock.RLock()
	handler := errorHandler
	lock.RUnlock()

	if handler == nil {
		if len(fns) > 0 {
			fns[0](w, err)
		} else if errcode.IsGrpcError(err) {
			// don't unwrap error and get status.Message(),
			// it hides the rpc error headers.
			http.Error(w, err.Error(), errcode.CodeFromGrpcError(err))
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		return
	}

	code, body := handler(err)
	if body == nil {
		w.WriteHeader(code)
		return
	}

	e, ok := body.(error)
	if ok {
		http.Error(w, e.Error(), code)
	} else {
		WriteJson(w, code, body)
	}
}

// Ok writes HTTP 200 OK into w.
func Ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// OkJson writes v into w with 200 OK.
func OkJson(w http.ResponseWriter, v interface{}) {
	WriteJson(w, http.StatusOK, v)
}

// SetErrorHandler sets the error handler, which is called on calling Error.
func SetErrorHandler(handler func(error) (int, interface{})) {
	lock.Lock()
	defer lock.Unlock()
	errorHandler = handler
}

// WriteJson writes v as json string into w with code.
func WriteJson(w http.ResponseWriter, code int, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(ContentType, header.JsonContentType)
	w.WriteHeader(code)

	if n, err := w.Write(bs); err != nil {
		// http.ErrHandlerTimeout has been handled by http.TimeoutHandler,
		// so it's ignored here.
		if err != http.ErrHandlerTimeout {
			logx.Errorf("write response failed, error: %s", err)
		}
	} else if n < len(bs) {
		logx.Errorf("actual bytes: %d, written bytes: %d", len(bs), n)
	}
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
	var body Body
	if err != nil {
		switch e := err.(type) {
		case *errorx.CodeError:
			body.Code = e.Code
			body.Msg = fmt.Sprintf(e.Msg, e.Placeholder...)
		default:
			body.Code = -1
			body.Msg = err.Error()
		}
	} else {
		body.Msg = "OK"
		body.Data = resp
	}
	OkJson(w, body)
}

func ResponseWithLang(r *http.Request, w http.ResponseWriter, l *lang.Lang, resp interface{}, err error) {
	var body Body
	if err != nil {
		switch e := err.(type) {
		case *errorx.CodeError:
			body.Code = e.Code
			langH := r.Header.Get("language")
			if len(langH) != 0 && len(e.LangCode) != 0 {
				body.Msg = l.ParseMsg(langH, e.LangCode, e.Placeholder...)
			}
		default:
			body.Code = -1
		}
		if len(body.Msg) == 0 {
			body.Msg = err.Error()
		}
	} else {
		body.Msg = "OK"
		body.Data = resp
	}
	OkJson(w, body)
}

func ResponseWithLangAndData(r *http.Request, w http.ResponseWriter, l *lang.Lang, resp interface{}, err error) {
	var body Body
	if err != nil {
		switch e := err.(type) {
		case *errorx.CodeError:
			body.Code = e.Code
			langH := r.Header.Get("language")
			if len(langH) != 0 && len(e.LangCode) != 0 {
				body.Msg = l.ParseMsg(langH, e.LangCode, e.Placeholder...)
			}
		default:
			body.Code = -1
		}
		if len(body.Msg) == 0 {
			body.Msg = err.Error()
		}
		body.Data = resp
	} else {
		body.Msg = "OK"
		body.Data = resp
	}
	OkJson(w, body)
}
