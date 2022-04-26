package api

import (
	"encoding/json"
	"gluwholevpp/pkg/repository"
	"gluwholevpp/pkg/vpp"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *Api) GetBitstreamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var bitstream repository.Bitstream

	err := a.storage.GetBitstream(vars["id"], &bitstream)
	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 500)
		return
	}

	res := &ResponseBitstream{}
	res.Status = 1
	res.Msg = &bitstream
	writeHttpResponseJSON(res, &w, 200)
}

func (a *Api) GetBitstreamsHandler(w http.ResponseWriter, r *http.Request) {
	var bitstreams []repository.Bitstream
	err := a.storage.GetBitstreams(&bitstreams)

	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 500)
		return
	}
	res := &ResponseBitstreams{}
	res.Status = 1
	res.Msg = &bitstreams
	writeHttpResponseJSON(res, &w, 200)
}

func (a *Api) GetBitstreamsFromCustomerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 400)
		return
	}

	var bitstreams []repository.Bitstream
	err = a.storage.GetBitstreamsFromCustomer(id, &bitstreams)

	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 500)
		return
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

	idCounter, err := a.storage.IncrementIdCounter()

	if err != nil {
		res.Msg = err.Error()
		writeHttpResponseJSON(res, &w, 500)
		return
	}

	// We set Outer Vlan based on S-VLAN Customer (OuterVlan)
	bitstream.DstOuter = customer.OuterVlan
	// We set Inner Vlan (C-VLAN) based on a customer counter
	bitstream.DstInner = counter

	bitstream.SrcId = idCounter
	bitstream.DstId = idCounter + 1

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

	a.vpp.CreateBitstream(vppBitstream, a.prio)

	res.Status = 1
	res.Msg = "Ok"
	writeHttpResponseJSON(res, &w, 200)
}

func (a *Api) DeleteBitstreamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var bitstream repository.Bitstream

	err := a.storage.GetBitstream(vars["id"], &bitstream)

	res := &ResponseGeneric{}
	res.Status = 0

	if err != nil {
		res.Msg = err.Error()
		writeHttpResponseJSON(res, &w, 400)
		return
	}

	vppBitstream := &vpp.Bitstream{
		SrcInterface: a.srcInterface,
		DstInterface: 0,
		SrcId:        bitstream.SrcId,
		DstId:        bitstream.DstId,
		SrcOuter:     bitstream.SrcOuter,
		SrcInner:     bitstream.SrcInner,
		DstOuter:     bitstream.DstOuter,
		DstInner:     bitstream.DstInner,
	}

	// Delete from DB
	err = a.storage.DeleteBitstream(bitstream.Id)
	if err != nil {
		res.Msg = err.Error()
		writeHttpResponseJSON(res, &w, 500)
		return
	}

	a.vpp.DeleteBitstream(vppBitstream)

	res.Status = 1
	res.Msg = "Ok"
	writeHttpResponseJSON(res, &w, 200)
}

func (a *Api) LoadBitstreamsStorage() error {
	var bitstreams []repository.Bitstream
	err := a.storage.GetBitstreams(&bitstreams)

	if err != nil {
		log.Println(err)
		return err
	}

	for _, v := range bitstreams {
		var customer repository.Customer
		err = a.storage.GetCustomer(v.CustomerId, &customer)

		if err != nil {
			log.Println(err)
			return err
		}

		bitstream := &vpp.Bitstream{
			SrcInterface: a.srcInterface,
			DstInterface: customer.OuterInterface,
			SrcId:        v.SrcId,
			DstId:        v.DstId,
			SrcOuter:     v.SrcOuter,
			SrcInner:     v.SrcInner,
			DstOuter:     v.DstOuter,
			DstInner:     v.DstInner,
		}
		a.vpp.CreateBitstream(bitstream, a.prio)
	}

	return nil
}
