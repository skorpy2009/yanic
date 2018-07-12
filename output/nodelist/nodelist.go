package nodelist

import (
	"chaos.expert/FreifunkBremen/yanic/lib/jsontime"
	"chaos.expert/FreifunkBremen/yanic/runtime"
)

// NodeList rewritten after: https://github.com/ffnord/ffmap-backend/blob/c33ebf62f013e18bf71b5a38bd058847340db6b7/lib/nodelist.py
type NodeList struct {
	Version   string        `json:"version"`
	Timestamp jsontime.Time `json:"updated_at"` // Timestamp of the generation
	List      []*Node       `json:"nodes"`
}

type Node struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Position *Position `json:"position,omitempty"`
	Status   struct {
		Online      bool          `json:"online"`
		LastContact jsontime.Time `json:"lastcontact"`
		Clients     uint32        `json:"clients"`
	} `json:"status"`
}

type Position struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func NewNode(n *runtime.Node) (node *Node) {
	if nodeinfo := n.Nodeinfo; nodeinfo != nil {
		node = &Node{
			ID:   nodeinfo.NodeID,
			Name: nodeinfo.Hostname,
		}
		if location := nodeinfo.Location; location != nil {
			node.Position = &Position{Lat: location.Latitude, Long: location.Longitude}
		}

		node.Status.Online = n.Online
		node.Status.LastContact = n.Lastseen
		if statistics := n.Statistics; statistics != nil && n.Online {
			node.Status.Clients = statistics.Clients.Total
		}
	}
	return
}

func transform(nodes *runtime.Nodes) *NodeList {
	nodelist := &NodeList{
		Version:   "1.0.1",
		Timestamp: jsontime.Now(),
	}

	for _, nodeOrigin := range nodes.List {
		node := NewNode(nodeOrigin)
		if node != nil {
			nodelist.List = append(nodelist.List, node)
		}
	}
	return nodelist
}
