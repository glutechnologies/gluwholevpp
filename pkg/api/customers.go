package api

import (
	"encoding/json"
	"gluwholevpp/pkg/repository"
	"net/http"
)

func (a *Api) GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
	var customers []repository.Customer
	err := a.storage.GetCustomers(&customers)

	if err != nil {
		resE := &ResponseGeneric{}
		resE.Status = 0
		resE.Msg = err.Error()
		writeHttpResponseJSON(resE, &w, 500)
	}
	res := &ResponseCustomers{}
	res.Status = 1
	res.Msg = &customers
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
