package router

import (
	"ac-gateway/help/fasthttputils"
	"ac-gateway/help/strutils"
	"ac-gateway/pkg/models/ddl"
	"github.com/valyala/fasthttp"
)

var charPosMap = make(map[rune]int)
var urlCharCount int

func init() {
	var urlCharArray = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~:/?#[]@!$&'()*+,;=%"
	for ix, char := range urlCharArray {
		charPosMap[char] = ix
	}

	urlCharCount = len(urlCharArray)
}

type RouterMatcher struct {
	hostRouterMap map[string]*RouterTrieTree
}

func (matcher *RouterMatcher) Match(ctx *fasthttp.RequestCtx) (*Upstream, bool) {
	host := fasthttputils.GetHeaderHost(ctx)
	tree, ok := matcher.hostRouterMap[host]
	if !ok {
		return nil, false
	}
	path := fasthttputils.GetRequestPath(ctx)
	return tree.SearchFirst(path)
}

type RouterTrieTree struct {
	root *RouterTreeNode
}
type RouterTreeNode struct {
	char     rune
	upstream *Upstream
	subNodes []*RouterTreeNode
}

func (upstream Upstream) GetNode() (string, bool) {
	return upstream.lb.ChooseByAddresses(upstream.nodeList)
}

func ParseUpstream(upstream *ddl.Upstream, options ...RouterTreeNodeOption) *Upstream {
	if upstream == nil {
		return nil
	}
	nodes := make(map[string]int)
	nodeList := make([]string, 0, len(upstream.Nodes))
	for k, v := range upstream.Nodes {
		nodes[k] = v
		for i := 0; i < v; i++ {
			nodeList = append(nodeList, k)
		}
	}
	u := &Upstream{
		typ:      upstream.Type,
		nodes:    nodes,
		nodeList: nodeList,
		lb:       NewRoundRobinLoadBalancer(nodeList),
	}
	if len(options) != 0 {
		for _, option := range options {
			option(u)
		}
	}
	return u
}

func BuildRouterMatcher(list []ddl.Router) *RouterMatcher {
	hostRouterMap := make(map[string]*RouterTrieTree, len(list)/2)
	for _, item := range list {
		if strutils.IsEmpty(item.Host) {
			continue
		}
		if hostRouterMap[item.Host] == nil {
			hostRouterMap[item.Host] = NewRouterTrieTree()
		}
		trieTree := hostRouterMap[item.Host]
		upstream := item.Upstream
		trieTree.PutString(item.Uri, &upstream)
	}

	return &RouterMatcher{
		hostRouterMap: hostRouterMap,
	}
}

// 创建空字典树
func NewRouterTrieTree() *RouterTrieTree {
	return &RouterTrieTree{
		root: newTreeNode(0, nil),
	}
}

// 搜索路径上遇到的第一个字符串
// uri: 请求路径
func (tree *RouterTrieTree) SearchFirst(uri string) (*Upstream, bool) {
	node := tree.root
	for _, char := range uri {
		node = node.findSubNode(char)
		if nil == node {
			return nil, false
		}
		if node.upstream != nil {
			return node.upstream, true
		}

	}

	return nil, false
}

func newTreeNode(char rune, upstream *ddl.Upstream, options ...RouterTreeNodeOption) *RouterTreeNode {
	node := &RouterTreeNode{
		char:     char,
		upstream: ParseUpstream(upstream, options...),
		subNodes: nil,
	}
	return node
}

// 添加一条path->serviceInfo映射
func (tree *RouterTrieTree) PutString(uri string, upstream *ddl.Upstream) {
	uriRunes := []rune(uri)
	LEN := len(uriRunes)

	node := tree.root
	for ix, char := range uriRunes {
		subNode := findNode(char, node.subNodes)
		if nil == subNode {
			var newNode *RouterTreeNode
			// 是最后一个字符
			if ix == LEN-1 {
				newNode = newTreeNode(char, upstream, ApplyRouterTreeNodeURI(uri))
			} else {
				newNode = newTreeNode(char, nil)
			}

			node.addSubNode(newNode)
			node = newNode

		} else if ix == LEN-1 {
			subNode.upstream = ParseUpstream(upstream, ApplyRouterTreeNodeURI(uri))
		} else {
			node = subNode
		}
	}
}

// 添加子节点
func (node *RouterTreeNode) addSubNode(newNode *RouterTreeNode) {
	if node.subNodes == nil {
		node.subNodes = make([]*RouterTreeNode, urlCharCount)
	}
	position := mapPosition(newNode.char)
	node.subNodes[position] = newNode
}
func (node *RouterTreeNode) findSubNode(target rune) *RouterTreeNode {
	if nil == node.subNodes {
		return nil
	}
	pos := mapPosition(target)
	return node.subNodes[pos]
}
func findNode(char rune, nodeList []*RouterTreeNode) *RouterTreeNode {
	if nil == nodeList {
		return nil
	}

	pos := mapPosition(char)
	return nodeList[pos]
}

func mapPosition(char rune) int {
	return charPosMap[char]
}
