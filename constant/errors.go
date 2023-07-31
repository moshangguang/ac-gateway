package constant

import "errors"

var (
	ErrCtxNotFoundRequestPath   = errors.New("err ctx not found request path")
	ErrCtxNotFoundRequestMethod = errors.New("err ctx not found request method")
	ErrCtxNotFoundMatchRouter   = errors.New("err ctx not found match router")
	ErrCtxNotFoundIP            = errors.New("err ctx not found ip")
	ErrRouterEmptyNodes         = errors.New("err router empty nodes")
	ErrRouterNodeSpace          = errors.New("err router node space")
	ErrRouterUpstreamType       = errors.New("err router upstream type")
	ErrRouterPluginEmpty        = errors.New("err router plugin empty")
	ErrRouterNotFoundPlugin     = errors.New("err router not found plugin")
	ErrHttpMethod               = errors.New("err http method")
	ErrRouterNotFound           = errors.New("err router not found")
	ErrRouterNameExists         = errors.New("err router name exists")
	ErrInitRouterFail           = errors.New("err init router fail")
)
