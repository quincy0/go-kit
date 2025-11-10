package router

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/quincy0/go-kit/core/search"
	"github.com/quincy0/go-kit/rest/httpx"
	"github.com/quincy0/go-kit/rest/pathvar"
)

const (
	allowHeader          = "Allow"
	allowMethodSeparator = ", "
)

var (
	// ErrInvalidMethod is an error that indicates not a valid http method.
	ErrInvalidMethod = errors.New("not a valid http method")
	// ErrInvalidPath is an error that indicates path is not start with /.
	ErrInvalidPath = errors.New("path must begin with '/'")
)

type patRouter struct {
	trees      map[string]*search.Tree
	notFound   http.Handler
	notAllowed http.Handler
}

// NewRouter returns a httpx.Router.
func NewRouter() httpx.Router {
	return &patRouter{
		trees: make(map[string]*search.Tree),
	}
}

func (pr *patRouter) Handle(method, reqPath string, handler http.Handler) error {
	if !validMethod(method) {
		return ErrInvalidMethod
	}

	if len(reqPath) == 0 || reqPath[0] != '/' {
		return ErrInvalidPath
	}

	cleanPath := path.Clean(reqPath)
	tree, ok := pr.trees[method]
	if ok {
		return tree.Add(cleanPath, handler)
	}

	tree = search.NewTree()
	pr.trees[method] = tree
	return tree.Add(cleanPath, handler)
}

func (pr *patRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ret := strings.Split(r.URL.Path, "/")
	tmpPath := ""
	if len(ret) > 1 {
		exts := strings.Split(ret[len(ret)-1], ".")
		if len(exts) > 1 && (exts[len(exts)-1] == "map" ||
			exts[len(exts)-1] == "json" ||
			exts[len(exts)-1] == "ico" ||
			exts[len(exts)-1] == "js" ||
			exts[len(exts)-1] == "html" ||
			exts[len(exts)-1] == "png" ||
			exts[len(exts)-1] == "jpg" ||
			exts[len(exts)-1] == "css") {
			tmpPath = r.URL.Path
			r.URL.Path = fmt.Sprintf("/%s/", ret[1])
		}
	}
	reqPath := path.Clean(r.URL.Path)
	if tree, ok := pr.trees[r.Method]; ok {
		if result, ok := tree.Search(reqPath); ok {
			if len(result.Params) > 0 {
				r = pathvar.WithVars(r, result.Params)
			}
			if len(tmpPath) > 0 {
				r.URL.Path = tmpPath
			}
			result.Item.(http.Handler).ServeHTTP(w, r)
			return
		}
	}

	allows, ok := pr.methodsAllowed(r.Method, reqPath)
	if !ok {
		pr.handleNotFound(w, r)
		return
	}

	if pr.notAllowed != nil {
		pr.notAllowed.ServeHTTP(w, r)
	} else {
		w.Header().Set(allowHeader, allows)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (pr *patRouter) SetNotFoundHandler(handler http.Handler) {
	pr.notFound = handler
}

func (pr *patRouter) SetNotAllowedHandler(handler http.Handler) {
	pr.notAllowed = handler
}

func (pr *patRouter) handleNotFound(w http.ResponseWriter, r *http.Request) {
	if pr.notFound != nil {
		pr.notFound.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (pr *patRouter) methodsAllowed(method, path string) (string, bool) {
	var allows []string

	for treeMethod, tree := range pr.trees {
		if treeMethod == method {
			continue
		}

		_, ok := tree.Search(path)
		if ok {
			allows = append(allows, treeMethod)
		}
	}

	if len(allows) > 0 {
		return strings.Join(allows, allowMethodSeparator), true
	}

	return "", false
}

func validMethod(method string) bool {
	return method == http.MethodDelete || method == http.MethodGet ||
		method == http.MethodHead || method == http.MethodOptions ||
		method == http.MethodPatch || method == http.MethodPost ||
		method == http.MethodPut
}
