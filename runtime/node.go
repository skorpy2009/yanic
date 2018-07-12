package runtime

import (
	"net"

	"chaos.expert/FreifunkBremen/yanic/data"
	"chaos.expert/FreifunkBremen/yanic/lib/jsontime"
)

// Node struct
type Node struct {
	Address    *net.UDPAddr     `json:"-"` // the last known address
	Firstseen  jsontime.Time    `json:"firstseen"`
	Lastseen   jsontime.Time    `json:"lastseen"`
	Online     bool             `json:"online"`
	Statistics *data.Statistics `json:"statistics"`
	Nodeinfo   *data.NodeInfo   `json:"nodeinfo"`
	Neighbours *data.Neighbours `json:"-"`
}

// Link represents a link between two nodes
type Link struct {
	SourceID      string
	SourceAddress string
	TargetID      string
	TargetAddress string
	TQ            float32
}

// IsGateway returns whether the node is a gateway
func (node *Node) IsGateway() bool {
	if info := node.Nodeinfo; info != nil {
		return info.VPN
	}
	return false
}
