package constant

import "net/http"

var DefaultHttpMethod = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodDelete,
	http.MethodHead,
	http.MethodPut,
	http.MethodPatch,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}
