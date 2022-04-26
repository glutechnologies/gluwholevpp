package api

import (
	"encoding/json"
	"gluwholevpp/pkg/repository"
	"net/http"
	"strconv"

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

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 400)
		return
	}

	var customer repository.Customer
	err = a.storage.GetCustomer(id, &customer)

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
