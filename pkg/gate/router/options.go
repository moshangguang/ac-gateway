package router

type RouterTreeNodeOption func(upstream *Upstream)

func ApplyRouterTreeNodeURI(uri string) RouterTreeNodeOption {
	return func(upstream *Upstream) {
		upstream.uri = uri
	}
}
