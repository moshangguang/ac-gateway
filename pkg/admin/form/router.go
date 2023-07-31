package form

import "ac-gateway/pkg/models/ddl"

type RouterAddForm struct {
	Host     string       `form:"host"`
	Name     string       `form:"name"`
	Uri      string       `form:"uri"`
	Methods  []string     `json:"methods"`
	Plugins  ddl.Plugins  `form:"plugins"`
	Priority int          `json:"priority"`
	Upstream ddl.Upstream `form:"upstream" binding:"required"`
}

func (form RouterAddForm) ToRouter() ddl.Router {
	return ddl.Router{
		Methods:  form.Methods,
		Name:     form.Name,
		Host:     form.Host,
		Uri:      form.Uri,
		Plugins:  form.Plugins,
		Priority: form.Priority,
		Upstream: form.Upstream,
	}
}

type RouterUpdateForm struct {
	Id       int64        `form:"host"`
	Name     string       `form:"name"`
	Host     string       `form:"host"`
	Uri      string       `form:"uri"`
	Methods  []string     `form:"methods"`
	Plugins  ddl.Plugins  `form:"plugins"`
	Priority int          `form:"priority"`
	Upstream ddl.Upstream `form:"upstream" binding:"required"`
}

func (form RouterUpdateForm) ToRouter() ddl.Router {
	return ddl.Router{
		Id:       form.Id,
		Name:     form.Name,
		Methods:  form.Methods,
		Host:     form.Host,
		Uri:      form.Uri,
		Plugins:  form.Plugins,
		Priority: form.Priority,
		Upstream: form.Upstream,
	}
}
