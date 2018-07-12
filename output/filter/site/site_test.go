package site

import (
	"testing"

	"chaos.expert/FreifunkBremen/yanic/data"
	"chaos.expert/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestFilterSite(t *testing.T) {
	assert := assert.New(t)

	// invalid config
	filter, err := build("ffhb")
	assert.Error(err)

	filter, err = build([]interface{}{3, "ffhb"})
	assert.Error(err)

	filter, err = build([]interface{}{"ffhb"})
	assert.NoError(err)

	// wronge node
	n := filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{System: data.System{SiteCode: "ffxx"}}})
	assert.Nil(n)

	// right node
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{System: data.System{SiteCode: "ffhb"}}})
	assert.NotNil(n)

	// node without data -> wrong node
	n = filter.Apply(&runtime.Node{})
	assert.Nil(n)
}
