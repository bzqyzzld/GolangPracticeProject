package gee

import "strings"

// 实现trie树结构,前缀树算法

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool //是否是模糊匹配, 含有 : 或者 * 的时候为true
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild == true {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild == true {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 注册路由的是由的使用, 递归使用
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		// 说明没有该part类型的规则,需要手动创建
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	// 递归插入
	child.insert(pattern, parts, height+1)

}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
