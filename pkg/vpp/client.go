package vpp

import (
	"log"

	"git.fd.io/govpp.git"
	interfaces "git.fd.io/govpp.git/binapi/interface"
	"git.fd.io/govpp.git/binapi/interface_types"
	"git.fd.io/govpp.git/binapi/l2"
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
	// If vpp is not enabled return
	if !c.enabled {
		return
	}

	c.conn.Disconnect()
}

func (c *Client) CreateBitstream(bitstream *Bitstream) {
	// If vpp is not enabled return
	if !c.enabled {
		return
	}

	// Create Src Vlan
	sw0, _ := CreateVlan(c.conn, bitstream.SrcInterface, bitstream.SrcId, bitstream.SrcOuter, bitstream.SrcInner)

	// Create Destination Vlan
	sw1, _ := CreateVlan(c.conn, bitstream.DstInterface, bitstream.DstId, bitstream.DstOuter, bitstream.DstInner)

	// Assign Vlan to the same bridge domain
	CreateBridgeDomain(c.conn, bitstream.SrcId, sw0, sw1)
}

func CreateVlan(c *core.Connection, sw int, id int, outer int, inner int) (interface_types.InterfaceIndex, error) {
	ch, err := c.NewAPIChannel()
	if err != nil {
		log.Println("ERROR: creating channel failed:", err)
		return 0, err
	}
	defer ch.Close()

	req := &interfaces.CreateSubif{SwIfIndex: interface_types.InterfaceIndex(sw), SubID: uint32(id), OuterVlanID: uint16(outer),
		InnerVlanID: uint16(inner), SubIfFlags: interface_types.SUB_IF_API_FLAG_TWO_TAGS}

	reply := &interfaces.CreateSubifReply{}

	if err = ch.SendRequest(req).ReceiveReply(reply); err != nil {
		log.Println("ERROR: creating sub-interface", err)
		return 0, err
	}

	req2 := &interfaces.SwInterfaceSetFlags{SwIfIndex: reply.SwIfIndex, Flags: interface_types.IF_STATUS_API_FLAG_ADMIN_UP}

	reply2 := &interfaces.SwInterfaceSetFlagsReply{}

	if err = ch.SendRequest(req2).ReceiveReply(reply2); err != nil {
		log.Println("ERROR: setting up interface", err)
		return 0, err
	}

	return reply.SwIfIndex, nil
}

func CreateBridgeDomain(c *core.Connection, id int, sw0 interface_types.InterfaceIndex, sw1 interface_types.InterfaceIndex) error {
	ch, err := c.NewAPIChannel()
	if err != nil {
		log.Println("ERROR: creating channel failed:", err)
		return err
	}
	defer ch.Close()

	req := &l2.BridgeDomainAddDel{
		BdID:    uint32(id),
		Learn:   true,
		Forward: true,
		Flood:   true,
		UuFlood: true,
		IsAdd:   true,
	}

	reply := &l2.BridgeDomainAddDelReply{}

	if err = ch.SendRequest(req).ReceiveReply(reply); err != nil {
		log.Println("ERROR: creating bridge-domain", err)
		return err
	}

	req2 := &l2.SwInterfaceSetL2Bridge{
		BdID:        uint32(id),
		RxSwIfIndex: interface_types.InterfaceIndex(sw0),
		Enable:      true,
	}

	reply2 := &l2.SwInterfaceSetL2BridgeReply{}

	if err = ch.SendRequest(req2).ReceiveReply(reply2); err != nil {
		log.Println("ERROR: seting sw0 to bridge-domain", err)
		return err
	}

	req3 := &l2.SwInterfaceSetL2Bridge{BdID: uint32(id),
		RxSwIfIndex: interface_types.InterfaceIndex(sw1),
		Enable:      true,
	}

	reply3 := &l2.SwInterfaceSetL2BridgeReply{}

	if err = ch.SendRequest(req3).ReceiveReply(reply3); err != nil {
		log.Println("ERROR: seting sw1 to bridge-domain", err)
		return err
	}

	// Pop 2 tags
	req4 := &l2.L2InterfaceVlanTagRewrite{SwIfIndex: interface_types.InterfaceIndex(sw0), VtrOp: 4}

	reply4 := &l2.L2InterfaceVlanTagRewriteReply{}

	if err = ch.SendRequest(req4).ReceiveReply(reply4); err != nil {
		log.Println("ERROR: seting tag rewrite pop2 sw0", err)
		return err
	}

	req5 := &l2.L2InterfaceVlanTagRewrite{SwIfIndex: interface_types.InterfaceIndex(sw1), VtrOp: 4}

	reply5 := &l2.L2InterfaceVlanTagRewriteReply{}

	if err = ch.SendRequest(req5).ReceiveReply(reply5); err != nil {
		log.Println("ERROR: seting tag rewrite pop2 sw1", err)
		return err
	}

	return nil
}
