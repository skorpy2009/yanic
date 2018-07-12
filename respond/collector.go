package respond

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"chaos.expert/FreifunkBremen/yanic/data"
	"chaos.expert/FreifunkBremen/yanic/database"
	"chaos.expert/FreifunkBremen/yanic/lib/jsontime"
	"chaos.expert/FreifunkBremen/yanic/runtime"
)

// Collector for a specificle respond messages
type Collector struct {
	connections []multicastConn // UDP sockets

	queue        chan *Response // received responses
	db           database.Connection
	nodes        *runtime.Nodes
	sitesDomains map[string][]string
	interval     time.Duration // Interval for multicast packets
	stop         chan interface{}
}

type multicastConn struct {
	Conn             *net.UDPConn
	SendRequest      bool
	MulticastAddress net.IP
}

// NewCollector creates a Collector struct
func NewCollector(db database.Connection, nodes *runtime.Nodes, sitesDomains map[string][]string, ifaces []InterfaceConfig) *Collector {

	coll := &Collector{
		db:           db,
		nodes:        nodes,
		sitesDomains: sitesDomains,
		queue:        make(chan *Response, 400),
		stop:         make(chan interface{}),
	}

	for _, iface := range ifaces {
		coll.listenUDP(iface)
	}

	go coll.parser()

	if coll.db != nil {
		go coll.globalStatsWorker()
	}

	return coll
}

func (coll *Collector) listenUDP(iface InterfaceConfig) {

	var addr net.IP

	var err error
	if iface.IPAddress != "" {
		addr = net.ParseIP(iface.IPAddress)
	} else {
		addr, err = getUnicastAddr(iface.InterfaceName, iface.MulticastAddress == "")
		if err != nil {
			log.Panic(err)
		}
	}

	multicastAddress := multicastAddressDefault
	if iface.MulticastAddress != "" {
		multicastAddress = iface.MulticastAddress
	}

	// Open socket
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   addr,
		Port: iface.Port,
		Zone: iface.InterfaceName,
	})
	if err != nil {
		log.Panic(err)
	}
	conn.SetReadBuffer(maxDataGramSize)

	coll.connections = append(coll.connections, multicastConn{
		Conn:             conn,
		SendRequest:      !iface.SendNoRequest,
		MulticastAddress: net.ParseIP(multicastAddress),
	})

	// Start receiver
	go coll.receiver(conn)
}

// Returns a unicast address of given interface (linklocal or global unicast address)
func getUnicastAddr(ifname string, linklocal bool) (net.IP, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}

	addresses, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	var ip net.IP

	for _, addr := range addresses {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if (!linklocal && ipnet.IP.IsGlobalUnicast()) || (linklocal && ipnet.IP.IsLinkLocalUnicast()) {
			ip = ipnet.IP
		}
	}
	if ip != nil {
		return ip, nil
	}
	return nil, fmt.Errorf("unable to find a unicast address for %s", ifname)
}

// Start Collector
func (coll *Collector) Start(interval time.Duration) {
	if coll.interval != 0 {
		panic("already started")
	}
	if interval <= 0 {
		panic("invalid collector interval")
	}
	coll.interval = interval

	go func() {
		coll.sendOnce() // immediately
		coll.sender()   // periodically
	}()
}

// Close Collector
func (coll *Collector) Close() {
	close(coll.stop)
	for _, conn := range coll.connections {
		conn.Conn.Close()
	}
	close(coll.queue)
}

func (coll *Collector) sendOnce() {
	now := jsontime.Now()
	coll.sendMulticast()

	// Wait for the multicast responses to be processed and send unicasts
	time.Sleep(coll.interval / 2)
	coll.sendUnicasts(now)
}

func (coll *Collector) sendMulticast() {
	log.Println("sending multicasts")
	for _, conn := range coll.connections {
		if conn.SendRequest {
			coll.sendPacket(conn.Conn, conn.MulticastAddress)
		}
	}
}

// Send unicast packets to nodes that did not answer the multicast
func (coll *Collector) sendUnicasts(seenBefore jsontime.Time) {
	seenAfter := seenBefore.Add(-time.Minute * 10)

	// Select online nodes that has not been seen recently
	nodes := coll.nodes.Select(func(n *runtime.Node) bool {
		return n.Lastseen.After(seenAfter) && n.Lastseen.Before(seenBefore) && n.Address != nil
	})

	// Send unicast packets
	count := 0
	for _, node := range nodes {
		send := 0
		for _, conn := range coll.connections {
			if node.Address.Zone != "" && conn.Conn.LocalAddr().(*net.UDPAddr).Zone != node.Address.Zone {
				continue
			}
			coll.sendPacket(conn.Conn, node.Address.IP)
			send++
		}
		if send == 0 {
			log.Printf("unable to find connection for %s", node.Address.Zone)
		} else {
			time.Sleep(10 * time.Millisecond)
			count += send
		}
	}
	log.Printf("sending %d unicast pkg for %d nodes", count, len(nodes))
}

// SendPacket sends a UDP request to the given unicast or multicast address on the first UDP socket
func (coll *Collector) SendPacket(destination net.IP) {
	coll.sendPacket(coll.connections[0].Conn, destination)
}

// sendPacket sends a UDP request to the given unicast or multicast address on the given UDP socket
func (coll *Collector) sendPacket(conn *net.UDPConn, destination net.IP) {
	addr := net.UDPAddr{
		IP:   destination,
		Port: port,
		Zone: conn.LocalAddr().(*net.UDPAddr).Zone,
	}

	if _, err := conn.WriteToUDP([]byte("GET nodeinfo statistics neighbours"), &addr); err != nil {
		log.Println("WriteToUDP failed:", err)
	}
}

// send packets continuously
func (coll *Collector) sender() {
	ticker := time.NewTicker(coll.interval)
	for {
		select {
		case <-coll.stop:
			ticker.Stop()
			return
		case <-ticker.C:
			// send the multicast packet to request per-node statistics
			coll.sendOnce()
		}
	}
}

func (coll *Collector) parser() {
	for obj := range coll.queue {
		if data, err := obj.parse(); err != nil {
			log.Println("unable to decode response from", obj.Address.String(), err)
		} else {
			coll.saveResponse(obj.Address, data)
		}
	}
}

func (res *Response) parse() (*data.ResponseData, error) {
	// Deflate
	deflater := flate.NewReader(bytes.NewReader(res.Raw))
	defer deflater.Close()

	// Unmarshal
	rdata := &data.ResponseData{}
	err := json.NewDecoder(deflater).Decode(rdata)

	return rdata, err
}

func (coll *Collector) saveResponse(addr *net.UDPAddr, res *data.ResponseData) {
	// Search for NodeID
	var nodeID string
	if val := res.NodeInfo; val != nil {
		nodeID = val.NodeID
	} else if val := res.Neighbours; val != nil {
		nodeID = val.NodeID
	} else if val := res.Statistics; val != nil {
		nodeID = val.NodeID
	}

	// Check length of nodeID
	if len(nodeID) != 12 {
		log.Printf("invalid NodeID '%s' from %s", nodeID, addr.String())
		return
	}

	// Set fields to nil if nodeID is inconsistent
	if res.Statistics != nil && res.Statistics.NodeID != nodeID {
		res.Statistics = nil
	}
	if res.Neighbours != nil && res.Neighbours.NodeID != nodeID {
		res.Neighbours = nil
	}
	if res.NodeInfo != nil && res.NodeInfo.NodeID != nodeID {
		res.NodeInfo = nil
	}

	// Process the data and update IP address
	node := coll.nodes.Update(nodeID, res)
	node.Address = addr

	// Store statistics in database
	if db := coll.db; db != nil {
		db.InsertNode(node)

		// Store link data
		if neighbours := node.Neighbours; neighbours != nil {
			coll.nodes.RLock()
			for _, link := range coll.nodes.NodeLinks(node) {
				db.InsertLink(&link, node.Lastseen.GetTime())
			}
			coll.nodes.RUnlock()
		}
	}
}

func (coll *Collector) receiver(conn *net.UDPConn) {
	buf := make([]byte, maxDataGramSize)
	for {
		n, src, err := conn.ReadFromUDP(buf)

		if err != nil {
			log.Println("ReadFromUDP failed:", err)
			return
		}

		raw := make([]byte, n)
		copy(raw, buf)

		coll.queue <- &Response{
			Address: src,
			Raw:     raw,
		}
	}
}

func (coll *Collector) globalStatsWorker() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-coll.stop:
			ticker.Stop()
			return
		case <-ticker.C:
			coll.saveGlobalStats()
		}
	}
}

// saves global statistics
func (coll *Collector) saveGlobalStats() {
	stats := runtime.NewGlobalStats(coll.nodes, coll.sitesDomains)

	for site, domains := range stats {
		for domain, stat := range domains {
			coll.db.InsertGlobals(stat, time.Now(), site, domain)
		}
	}
}
