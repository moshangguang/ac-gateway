package library

import (
	"ac-gateway/constant"
	"ac-gateway/help/strutils"
	"ac-gateway/pkg/models/ddl"
	"github.com/wxnacy/wgo/arrays"
	"strings"
)

func FormatRouter(data ddl.Router) ddl.Router {
	uri := data.Uri
	if strutils.IsEmpty(data.Uri) {
		uri = "/**"
	}
	methods := make([]string, 0, len(constant.DefaultHttpMethod))
	if len(data.Methods) == 0 {
		methods = constant.DefaultHttpMethod
	} else {
		for _, m := range data.Methods {
			m = strings.ToUpper(m)
			if arrays.ContainsString(constant.DefaultHttpMethod, m) == -1 {
				continue
			}
			if arrays.ContainsString(methods, m) == -1 {
				methods = append(methods, m)
			}
		}
	}
	return ddl.Router{
		Id:       data.Id,
		Methods:  methods,
		Host:     data.Host,
		Uri:      uri,
		Plugins:  data.Plugins,
		Upstream: data.Upstream,
	}
}
func FormatRouterList(data []ddl.Router) []ddl.Router {
	result := make([]ddl.Router, 0, len(data))
	for _, d := range data {
		result = append(result, FormatRouter(d))
	}
	return result
}
