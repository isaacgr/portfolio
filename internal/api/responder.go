package api

import (
	"net/http"

	"github.com/segmentio/encoding/json"
)

type Responder interface {
	JSON(w http.ResponseWriter, status int, data any)
	HTML(w http.ResponseWriter, status int, data string)
	Error(w http.ResponseWriter, status int, err error)
	DefineResponseHeaders(w http.ResponseWriter)
}

type responder struct{}

func NewResponder() *responder {
	return &responder{}
}

func (r *responder) DefineResponseHeaders(w http.ResponseWriter) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	w.Header().Set(
		"Access-Control-Allow-Origin",
		"*",
	)
	w.Header().Set(
		"Access-Control-Allow-Methods",
		"GET, PATCH, PUT, POST, HEAD, OPTIONS, DELETE",
	)
	w.Header().Set(
		"Access-Control-Allow-Headers",
		"Content-Type, Accept",
	)
	w.Header().Set(
		"Access-Control-Max-Age",
		"3600",
	)
}

func (r *responder) JSON(w http.ResponseWriter, status int, data any) {
	r.DefineResponseHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(data)
	if err != nil {
		r.Error(
			w,
			http.StatusInternalServerError,
			Error{
				Msg:  "Failed to encode response",
				Code: http.StatusInternalServerError,
				Data: "",
			},
		)
		return
	}
	w.WriteHeader(status)
	w.Write(res)
}

func (r *responder) HTML(
	w http.ResponseWriter,
	status int,
	data string,
) {
	r.DefineResponseHeaders(w)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	w.Write([]byte(data))
}

func (r *responder) Error(
	w http.ResponseWriter,
	status int,
	err error,
) {
	if e, ok := err.(Error); ok {
		r.JSON(w, status, e)
	} else {
		e := Error{
			Msg:  err.Error(),
			Code: http.StatusInternalServerError,
			Data: "",
		}
		r.JSON(w, status, e)
	}
}
