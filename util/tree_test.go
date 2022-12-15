package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkNewTree(b *testing.B) {
	tree := NewRootTree()
	tree.AddNode("/test/a/b/c", "1")
	tree.AddNode("/a/b/c/d", "2")
	tree.AddNode("/a/b/.+", "3")
	tree.AddNode("/c/.+", "4")
	tree.AddNode("/3/.+", "5")
	tree.AddNode("/2/.+", "6")
	tree.AddNode("/1/.+", "7")
	for i := 0; i < b.N; i++ {
		tree.IsMatch("/test/a/b/c")
	}
}

func TestTreeAdd(t *testing.T) {
	tree := NewRootTree()
	tree.AddNode("/test/a/b/c", "1")
	tree.AddNode("/a/b/c/d", "2")
	tree.AddNode("/a/b/.+", "3")
	v, serviceId := tree.IsMatch("/test/a/b/c")
	assert.Equal(t, true, v)
	assert.Equal(t, "1", serviceId.serviceID)
	v, serviceId = tree.IsMatch("/a/b/c")
	assert.Equal(t, true, v)
	assert.Equal(t, "3", serviceId.serviceID)
	v, serviceId = tree.IsMatch("/a/b/c/d")
	assert.Equal(t, true, v)
	assert.Equal(t, "2", serviceId.serviceID)
	tree.DelNode("/a/b/c/d")
	v, serviceId = tree.IsMatch("/a/b/c/d")
	assert.Equal(t, true, v)
	assert.Equal(t, "3", serviceId.serviceID)
	tree.AddNode("/a/b/c/d", "4")
	v, serviceId = tree.IsMatch("/a/b/c/d")
	assert.Equal(t, true, v)
	assert.Equal(t, "4", serviceId.serviceID)
}

func TestTreeNode(t *testing.T) {
	tree := NewRootTree()
	//tree.AddNode("/test/a/b/c", "1")
	tree.AddNode("/.+", "2")
	v, serviceId := tree.IsMatch("/a")
	fmt.Println(v, serviceId)
}
