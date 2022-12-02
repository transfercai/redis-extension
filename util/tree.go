package util

import (
	"strings"
)

var nodeTree = &Tree{
	isRoot: true,
	route:  NewSafeMap(),
}

type Tree struct {
	father    *Tree
	route     *SafeMap
	value     string
	isRoot    bool
	serviceID string
}

func NewTree() *Tree {
	return &Tree{
		route: NewSafeMap(),
	}
}

func (t *Tree) newChild(node string) *Tree {
	return &Tree{
		father: t,
		value:  node,
		route:  NewSafeMap(),
	}
}

func (t *Tree) root() *Tree {
	return nodeTree
}

func (t *Tree) AddNode(path, serviceID string) {
	nodes := strings.Split(path, "/")
	if len(nodes) == 0 {
		return
	}
	if nodes[0] == "" {
		nodes = nodes[1:]
	}
	n := t.root()
	count := len(nodes)
	for index, node := range nodes {
		if _, ok := n.route.Get(node); !ok {
			n.route.Set(node, n.newChild(node))
		}
		v, _ := n.route.Get(node)
		n = v.(*Tree)
		if index == count-1 {
			n.serviceID = serviceID
		}
	}
}

func (t *Tree) DelNode(path string) {
	ok, node := t.IsMatch(path)
	if !ok {
		return
	}
	node.serviceID = ""
	for !node.father.isRoot {
		if node.father.route.Len() < 2 {
			node.father.route.Del(node.value)
			tmp := node.father
			node.father = nil
			node = tmp
			continue
		}
		node.father.route.Del(node.value)
		break
	}
}

func (t *Tree) IsMatch(path string) (bool, *Tree) {
	nodes := strings.Split(path, "/")
	if len(nodes) == 0 {
		return false, nil
	}
	if nodes[0] == "" {
		nodes = nodes[1:]
	}
	n := t.root()
	count, index := len(nodes), 0
	for index < count {
		if v, ok := n.route.Get(nodes[index]); ok {
			n = v.(*Tree)
			index++
			continue
		}
		break
	}
	if index == count && n.serviceID != "" {
		return true, n
	}
	for !n.father.isRoot {
		if v, ok := n.father.route.Get(".+"); ok {
			return true, v.(*Tree)
		}
		n = n.father
	}
	return false, nil
}
