package blacklist

import (
	"testing"

	"chaos.expert/FreifunkBremen/yanic/data"
	"chaos.expert/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestFilterBlacklist(t *testing.T) {
	assert := assert.New(t)

	// invalid config
	filter, err := build(3)
	assert.Error(err)

	filter, err = build([]interface{}{2, "a"})
	assert.Error(err)

	// tests with empty list
	filter, err = build([]interface{}{})
	assert.NoError(err)

	// keep node without nodeid
	n := filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{}})
	assert.NotNil(n)

	// tests with blacklist
	filter, err = build([]interface{}{"a", "c"})
	assert.NoError(err)

	// blacklist contains node with nodeid -> drop it
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{NodeID: "a"}})
	assert.Nil(n)

	// blacklist does not contains node without nodeid -> keep it
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{}})
	assert.NotNil(n)

	n = filter.Apply(&runtime.Node{})
	assert.NotNil(n)
}
