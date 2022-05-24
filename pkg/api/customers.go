package api

import (
	"encoding/json"
	"gluwholevpp/pkg/repository"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// swagger:route GET /customers/ customer getCustomers
// Get a list of all customers
//
// responses:
//  401: ResponseGeneric
//	500: ResponseGeneric
//  200: ResponseCustomers
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

// swagger:route GET /customer/{id} customer getCustomer
// Get customer information
//
// responses:
//  401: ResponseGeneric
//	500: ResponseGeneric
//  200: ResponseCustomer
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

// swagger:route DELETE /customer/{id} customer deleteCustomer
// Delete given customer by id
//
// responses:
//  401: ResponseGeneric
//	500: ResponseGeneric
//  200: ResponseGeneric
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

// swagger:route POST /customer customer createCustomer
//  Create a customer
//
// parameters:
//  + name: id
//      in: body
//      description: Unique identifier for customer
//      required: true
//      type: string
//	+ name: name
//      in: body
//      description: Name for customer
//      required: true
//      type: string
//	+ name: outer-interface
//      in: body
//      description: Outer interface id from VPP Dataplane
//      required: true
//      type: int
//  + name: outer-vlan
//      in: body
//      description: Outer destination VLAN for customer (S-VLAN) where all customer's bitstreams will be included
//      required: true
//      type: int
//  + name: counter
//      in: body
//      description: C-VLAN counter incremented for each new bitstream
//      required: true
//      type: int
//
// responses:
//  401: ResponseGeneric
//	500: ResponseGeneric
//  200: ResponseGeneric
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

// swagger:route PATCH /customer/{id} customer createCustomer
//  Patch given customer with new values from body
//
// parameters:
//	+ name: name
//      in: body
//      description: Name for customer
//      required: true
//      type: string
//	+ name: outer-interface
//      in: body
//      description: Outer interface id from VPP Dataplane
//      required: true
//      type: int
//  + name: outer-vlan
//      in: body
//      description: Outer destination VLAN for customer (S-VLAN) where all customer's bitstreams will be included
//      required: true
//      type: int
//  + name: counter
//      in: body
//      description: C-VLAN counter incremented for each new bitstream
//      required: true
//      type: int
// responses:
//  401: ResponseGeneric
//	500: ResponseGeneric
//  200: ResponseGeneric
func (a *Api) PatchCustomerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	var customer repository.Customer
	err := a.storage.GetCustomer(id, &customer)

	res := &ResponseGeneric{}
	res.Status = 0

	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 400)
		return
	}

	// Decode body
	err = json.NewDecoder(r.Body).Decode(&customer)

	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 400)
		return
	}

	// Invalidate Id from decoder
	customer.Id = id

	// v10 validator for structs
	validate := validator.New()
	err = validate.Struct(customer)

	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 400)
		return
	}

	err = a.storage.UpdateCustomer(&customer)
	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 500)
		return
	}

	res.Status = 1
	res.Msg = "Ok"
	writeHttpResponseJSON(res, &w, 200)
}
