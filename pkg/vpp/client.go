package vpp

import (
	"gluwholevpp/pkg/utils"
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

func (c *Client) CreateBitstream(bitstream *Bitstream, prio int) {
	// If vpp is not enabled return
	if !c.enabled {
		return
	}

	// Create Src Vlan
	sw0, _ := CreateVlan(c.conn, bitstream.SrcInterface, bitstream.SrcId, bitstream.SrcOuter, bitstream.SrcInner)

	// Create Destination Vlan
	sw1, _ := CreateVlan(c.conn, bitstream.DstInterface, bitstream.DstId, bitstream.DstOuter, bitstream.DstInner)

	// Assign Vlan to the same bridge domain

	outer := utils.ConcatVlanPrio(bitstream.SrcOuter, prio)
	inner := utils.ConcatVlanPrio(bitstream.SrcInner, prio)

	CreateBridgeDomain(c.conn, bitstream.SrcId, sw0, sw1, outer, inner)
}

func (c *Client) DeleteBitstream(bitstream *Bitstream) {
	// If vpp is not enabled return
	if !c.enabled {
		return
	}

	// Prepare sub-interface dump
	ifaces := make(map[int]*interfaces.SwInterfaceDetails)
	SubInterfaceDump(c.conn, ifaces)

	// Delete Source VLAN
	DeleteVlan(c.conn, int(ifaces[bitstream.SrcId].SwIfIndex))

	// Delete Dest VLAN
	DeleteVlan(c.conn, int(ifaces[bitstream.DstId].SwIfIndex))

	// Delete bridge domain
	DeleteBridgeDomain(c.conn, bitstream.SrcId)
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

func CreateBridgeDomain(c *core.Connection, id int, sw0 interface_types.InterfaceIndex,
	sw1 interface_types.InterfaceIndex, outer int, inner int) error {
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
	req4 := &l2.L2InterfaceVlanTagRewrite{SwIfIndex: interface_types.InterfaceIndex(sw1), VtrOp: 8,
		Tag1: uint32(outer), Tag2: uint32(inner), PushDot1q: 1}

	reply4 := &l2.L2InterfaceVlanTagRewriteReply{}

	if err = ch.SendRequest(req4).ReceiveReply(reply4); err != nil {
		log.Println("ERROR: seting tag rewrite translate 2:2", err)
		return err
	}

	return nil
}

func DeleteVlan(c *core.Connection, sw int) error {
	ch, err := c.NewAPIChannel()
	if err != nil {
		log.Println("ERROR: creating channel failed:", err)
		return err
	}
	defer ch.Close()

	req := &interfaces.DeleteSubif{SwIfIndex: interface_types.InterfaceIndex(sw)}

	reply := &interfaces.DeleteSubifReply{}

	if err = ch.SendRequest(req).ReceiveReply(reply); err != nil {
		log.Println("ERROR: deleting sub-interface", err)
		return err
	}

	return nil
}

func DeleteBridgeDomain(c *core.Connection, id int) error {
	ch, err := c.NewAPIChannel()
	if err != nil {
		log.Println("ERROR: creating channel failed:", err)
		return err
	}
	defer ch.Close()

	req := &l2.BridgeDomainAddDel{
		BdID:  uint32(id),
		IsAdd: false,
	}

	reply := &l2.BridgeDomainAddDelReply{}

	if err = ch.SendRequest(req).ReceiveReply(reply); err != nil {
		log.Println("ERROR: deleting bridge-domain", err)
		return err
	}

	return nil
}

func SubInterfaceDump(c *core.Connection, ifaces map[int]*interfaces.SwInterfaceDetails) error {
	ch, err := c.NewAPIChannel()
	if err != nil {
		log.Println("ERROR: creating channel failed:", err)
		return err
	}
	defer ch.Close()

	reqCtx := ch.SendMultiRequest(&interfaces.SwInterfaceDump{
		SwIfIndex: ^interface_types.InterfaceIndex(0),
	})
	for {
		msg := &interfaces.SwInterfaceDetails{}
		stop, err := reqCtx.ReceiveReply(msg)
		if stop {
			break
		}
		if err != nil {
			log.Println(err, "dumping interfaces")
			return err
		}

		// Take only sub-interfaces
		if msg.SubID != 0 {
			ifaces[int(msg.SubID)] = msg
		}
	}

	return nil
}
