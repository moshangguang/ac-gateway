package ddl

type Router struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name,omitempty"`
	Methods  []string `json:"methods,omitempty"`
	Host     string   `json:"host,omitempty"`
	Uri      string   `json:"uri,omitempty"`
	Priority int      `json:"priority,omitempty"`
	Plugins  Plugins  `json:"plugins,omitempty"`
	Upstream Upstream `json:"upstream" form:"upstream" binding:"required"`
}
type Upstream struct {
	Type  string         `json:"type,omitempty" form:"type" binding:"required"`
	Nodes map[string]int `json:"nodes" form:"nodes" binding:"required"`
}

type Plugins struct {
	ExtPluginPreReq   ExtPluginPreReq   `json:"ext-plugin-pre-req,omitempty" form:"ext-plugin-pre-req"`
	ExtPluginPostResp ExtPluginPostResp `json:"ext-plugin-post-resp,omitempty" form:"ext-plugin-post-resp"`
}

type ExtPluginPreReq struct {
	Conf []ExtPluginConf `json:"conf,omitempty" form:"conf"`
}
type ExtPluginConf struct {
	Name  string `json:"name,omitempty" form:"name"`
	Value string `json:"value,omitempty" form:"value"`
}
type ExtPluginPostResp struct {
	Conf []ExtPluginConf `json:"conf,omitempty"`
}
