package api

import (
	"encoding/json"
	"gluwholevpp/pkg/repository"
	"gluwholevpp/pkg/utils"
	"gluwholevpp/pkg/vpp"
	"net/http"
)

func (a *Api) GetBitstreamsHandler(w http.ResponseWriter, r *http.Request) {
	var bitstreams []repository.Bitstream
	err := a.storage.GetBitstreams(&bitstreams)

	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 500)
	}
	res := &ResponseBitstreams{}
	res.Status = 1
	res.Msg = &bitstreams
	writeHttpResponseJSON(res, &w, 200)
}

func (a *Api) CreateBitstreamHandler(w http.ResponseWriter, r *http.Request) {
	var bitstream repository.Bitstream
	err := json.NewDecoder(r.Body).Decode(&bitstream)

	res := &ResponseGeneric{}
	res.Status = 0

	if err != nil {
		res.Msg = err.Error()
		writeHttpResponseJSON(res, &w, 400)
		return
	}

	var customer repository.Customer
	err = a.storage.GetCustomer(bitstream.CustomerId, &customer)

	if err != nil {
		res.Msg = err.Error()
		writeHttpResponseJSON(res, &w, 400)
		return
	}
	counter, err := a.storage.IncrementCounterCustomer(bitstream.CustomerId)

	if err != nil {
		res.Msg = err.Error()
		writeHttpResponseJSON(res, &w, 500)
		return
	}

	// We set Outer Vlan based on S-VLAN Customer (OuterVlan)
	bitstream.DstOuter = customer.OuterVlan
	// We set Inner Vlan (C-VLAN) based on a customer counter
	bitstream.DstInner = counter

	bitstream.SrcId = utils.GetSubInterfaceId(bitstream.CustomerId, bitstream.SrcOuter, bitstream.SrcInner)
	bitstream.DstId = utils.GetSubInterfaceId(bitstream.CustomerId, bitstream.DstOuter, bitstream.DstInner)

	err = a.storage.InsertBitstream(&bitstream)
	if err != nil {
		res.Msg = err.Error()
		writeHttpResponseJSON(res, &w, 500)
		return
	}

	vppBitstream := &vpp.Bitstream{
		SrcInterface: a.srcInterface,
		DstInterface: customer.OuterInterface,
		SrcId:        bitstream.SrcId,
		DstId:        bitstream.DstId,
		SrcOuter:     bitstream.SrcOuter,
		SrcInner:     bitstream.SrcInner,
		DstOuter:     bitstream.DstOuter,
		DstInner:     bitstream.DstInner,
	}

	a.vpp.CreateBitstream(vppBitstream)

	res.Status = 1
	res.Msg = "Ok"
	writeHttpResponseJSON(res, &w, 200)
}
