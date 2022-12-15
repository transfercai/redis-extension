package context

import (
	"fmt"
	"sync"

	"github.com/transfercai/redis-extension/util"
)

var r *Route

type Route struct {
	host sync.Map
}

func GetRoute() *Route {
	r = &Route{}
	return r
}

func (r *Route) IsMatch(host, path string) (bool, string) {
	v, ok := r.host.Load(host)
	if !ok {
		fmt.Println("invalid host")
		return false, ""
	}
	t := v.(*util.Tree)
	isMatch, tree := t.IsMatch(path)
	if !isMatch {
		return false, ""
	}
	return isMatch, tree.GetServiceID()
}

func (r *Route) UpsertHost(host string, tree *util.Tree) {
	r.host.Store(host, tree)
}
