package noowner

import (
	"errors"

	"chaos.expert/FreifunkBremen/yanic/data"
	"chaos.expert/FreifunkBremen/yanic/output/filter"
	"chaos.expert/FreifunkBremen/yanic/runtime"
)

type noowner struct{ has bool }

func init() {
	filter.Register("no_owner", build)
}

func build(config interface{}) (filter.Filter, error) {
	if value, ok := config.(bool); ok {
		return &noowner{has: value}, nil
	}
	return nil, errors.New("invalid configuration, boolean expected")
}

func (no *noowner) Apply(node *runtime.Node) *runtime.Node {
	if nodeinfo := node.Nodeinfo; nodeinfo != nil && no.has {
		node = &runtime.Node{
			Address:    node.Address,
			Firstseen:  node.Firstseen,
			Lastseen:   node.Lastseen,
			Online:     node.Online,
			Statistics: node.Statistics,
			Nodeinfo: &data.NodeInfo{
				NodeID:   nodeinfo.NodeID,
				Network:  nodeinfo.Network,
				System:   nodeinfo.System,
				Owner:    nil,
				Hostname: nodeinfo.Hostname,
				Location: nodeinfo.Location,
				Software: nodeinfo.Software,
				Hardware: nodeinfo.Hardware,
				VPN:      nodeinfo.VPN,
				Wireless: nodeinfo.Wireless,
			},
			Neighbours: node.Neighbours,
		}
	}
	return node
}
