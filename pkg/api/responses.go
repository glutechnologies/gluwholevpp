package api

import (
	"encoding/json"
	"gluwholevpp/pkg/repository"
	"net/http"
)

type Response interface {
	GetMsg() string
	GetStatus() int
}

// swagger:model ResponseGeneric
type ResponseGeneric struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func NewResponseGeneric(status int, msg string) *ResponseGeneric {
	res := &ResponseGeneric{
		Status: status,
		Msg:    msg,
	}
	return res
}

func (r *ResponseGeneric) GetMsg() string {
	return r.Msg
}

func (r *ResponseGeneric) GetStatus() int {
	return r.Status
}

// swagger:model ResponseBitstreams
type ResponseBitstreams struct {
	Status int                     `json:"status"`
	Msg    *[]repository.Bitstream `json:"msg"`
}

func (r *ResponseBitstreams) GetMsg() string {
	res, _ := json.Marshal(r.Msg)
	return string(res)
}

func (r *ResponseBitstreams) GetStatus() int {
	return r.Status
}

// swagger:model ResponseBitstream
type ResponseBitstream struct {
	Status int                   `json:"status"`
	Msg    *repository.Bitstream `json:"msg"`
}

func (r *ResponseBitstream) GetMsg() string {
	res, _ := json.Marshal(r.Msg)
	return string(res)
}

func (r *ResponseBitstream) GetStatus() int {
	return r.Status
}

// swagger:model ResponseCustomers
type ResponseCustomers struct {
	Status int                    `json:"status"`
	Msg    *[]repository.Customer `json:"msg"`
}

func (r *ResponseCustomers) GetMsg() string {
	res, _ := json.Marshal(r.Msg)
	return string(res)
}

func (r *ResponseCustomers) GetStatus() int {
	return r.Status
}

// swagger:model ResponseCustomer
type ResponseCustomer struct {
	Status int                  `json:"status"`
	Msg    *repository.Customer `json:"msg"`
}

func (r *ResponseCustomer) GetMsg() string {
	res, _ := json.Marshal(r.Msg)
	return string(res)
}

func (r *ResponseCustomer) GetStatus() int {
	return r.Status
}

func writeHttpResponseJSON(res Response, w *http.ResponseWriter, httpCode int) {
	r := &ResponseGeneric{}
	r.Status = res.GetStatus()
	r.Msg = res.GetMsg()
	json.NewEncoder((*w)).Encode(res)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(httpCode)
}
