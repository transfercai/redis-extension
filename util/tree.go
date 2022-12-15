package util

import (
	"strings"
)

type Tree struct {
	father    *Tree
	route     *SafeMap
	value     string
	isRoot    bool
	root      *Tree
	serviceID string
}

func NewRootTree() *Tree {
	t := &Tree{
		father:    nil,
		route:     NewSafeMap(),
		isRoot:    true,
		value:     "/",
		serviceID: "",
	}
	t.root = t
	return t
}

func (t *Tree) newChild(node string) *Tree {
	return &Tree{
		father: t,
		value:  node,
		root:   t.root,
		isRoot: false,
		route:  NewSafeMap(),
	}
}

func (t *Tree) getRoot() *Tree {
	return t.root
}

func (t *Tree) GetServiceID() string {
	return t.serviceID
}

func (t *Tree) AddNode(path, serviceID string) {
	nodes := strings.Split(path, "/")
	if len(nodes) == 0 {
		return
	}
	if nodes[0] == "" {
		nodes = nodes[1:]
	}
	n := t.getRoot()
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
	n := t.getRoot()
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
	for n != nil {
		if v, ok := n.route.Get(".+"); ok {
			return true, v.(*Tree)
		}
		n = n.father
	}
	return false, nil
}
