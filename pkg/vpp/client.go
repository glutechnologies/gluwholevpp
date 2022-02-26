package vpp

import (
	"log"

	"git.fd.io/govpp.git"
	interfaces "git.fd.io/govpp.git/binapi/interface"
	"git.fd.io/govpp.git/binapi/interface_types"
	"git.fd.io/govpp.git/core"
)

type Bitstream struct {
	SrcInterface int
	DstInterface int
	SrcId        int
	DstId        int
	SrcOuter     int
	SrcInner     int
	DstOuter     int
	DstInner     int
}

type Client struct {
	sockAddr string
	conn     *core.Connection
	enabled  bool
}

func (c *Client) Init(sockAddr string, enabled bool) {
	// Initialize all struct members
	c.sockAddr = sockAddr
	c.enabled = enabled

	// If vpp is not enabled return
	if !c.enabled {
		return
	}

	conn, connEv, err := govpp.AsyncConnect(sockAddr, core.DefaultMaxReconnectAttempts, core.DefaultReconnectInterval)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	c.conn = conn
	// wait for Connected event
	e := <-connEv
	if e.State != core.Connected {
		log.Fatalln("ERROR: connecting to VPP failed:", e.Error)
	}

}

func (c *Client) Close() {
	c.conn.Disconnect()
}

func (c *Client) CreateBitstream(bitstream *Bitstream) {
	// If vpp is not enabled return
	if !c.enabled {
		return
	}

	// Create Src Vlan
	CreateVlan(c.conn, bitstream.SrcInterface, bitstream.SrcId, bitstream.SrcOuter, bitstream.SrcInner)

	// Create Destination Vlan
	CreateVlan(c.conn, bitstream.DstInterface, bitstream.DstId, bitstream.DstOuter, bitstream.DstInner)
}

func CreateVlan(c *core.Connection, sw int, id int, outer int, inner int) {
	ch, err := c.NewAPIChannel()
	if err != nil {
		log.Fatalln("ERROR: creating channel failed:", err)
	}
	defer ch.Close()

	req := &interfaces.CreateSubif{SwIfIndex: interface_types.InterfaceIndex(sw), SubID: uint32(id), OuterVlanID: uint16(outer),
		InnerVlanID: uint16(inner), SubIfFlags: interface_types.SUB_IF_API_FLAG_EXACT_MATCH}

	reply := &interfaces.CreateSubifReply{}

	if err = ch.SendRequest(req).ReceiveReply(reply); err != nil {
		log.Fatalln("ERROR: creating sub-interface", err)
	}

	req2 := &interfaces.SwInterfaceSetFlags{SwIfIndex: reply.SwIfIndex, Flags: interface_types.IF_STATUS_API_FLAG_ADMIN_UP}

	reply2 := &interfaces.SwInterfaceSetFlagsReply{}

	if err = ch.SendRequest(req2).ReceiveReply(reply2); err != nil {
		log.Fatalln("ERROR: setting up interface", err)
	}
}

func CreateBridgeDomain(c *core.Connection, id int, sw0 int, sw1 int) {
	ch, err := c.NewAPIChannel()
	if err != nil {
		log.Fatalln("ERROR: creating channel failed:", err)
	}
	defer ch.Close()
}
