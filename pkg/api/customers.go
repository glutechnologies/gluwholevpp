package api

import (
	"encoding/json"
	"gluwholevpp/pkg/repository"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func (a *Api) GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
	var customers []repository.Customer
	err := a.storage.GetCustomers(&customers)

	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 500)
		return
	}
	res := &ResponseCustomers{}
	res.Status = 1
	res.Msg = &customers
	writeHttpResponseJSON(res, &w, 200)
}

func (a *Api) GetCustomerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	var customer repository.Customer
	err := a.storage.GetCustomer(id, &customer)

	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 400)
		return
	}
	res := &ResponseCustomer{}
	res.Status = 1
	res.Msg = &customer
	writeHttpResponseJSON(res, &w, 200)
}

func (a *Api) DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	var bitstreams []repository.Bitstream
	err := a.storage.GetBitstreamsFromCustomer(id, &bitstreams)

	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 400)
		return
	}

	if len(bitstreams) > 0 {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = "Existing bitstreams for this customer"
		writeHttpResponseJSON(resE, &w, 400)
		return
	}

	err = a.storage.DeleteCustomer(id)

	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 500)
		return
	}

	res := &ResponseGeneric{}
	res.Status = 1
	res.Msg = "Ok"
	writeHttpResponseJSON(res, &w, 200)
}

func (a *Api) CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	var customer repository.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)

	res := &ResponseGeneric{}
	res.Status = 0

	if err != nil {
		res.Msg = err.Error()
		writeHttpResponseJSON(res, &w, 400)
		return
	}

	// v10 validator for structs
	validate := validator.New()
	err = validate.Struct(customer)

	if err != nil {
		res.Msg = err.Error()
		writeHttpResponseJSON(res, &w, 400)
		return
	}

	err = a.storage.InsertCustomer(&customer)

	if err != nil {
		res.Msg = err.Error()
		writeHttpResponseJSON(res, &w, 500)
		return
	}

	res.Status = 1
	res.Msg = "ok"
	writeHttpResponseJSON(res, &w, 200)
}
