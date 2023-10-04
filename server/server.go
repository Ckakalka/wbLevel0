package server

import (
	"fmt"

	"net/http"

	"github.com/Ckakalka/wbLevel0/models"
)

type Http struct {
	http.Server
	orderCash *models.OrderCash
}

func NewHttp(addr string, cash *models.OrderCash) *Http {
	var server Http
	server.Addr = addr
	server.Handler = nil
	server.orderCash = cash
	http.DefaultServeMux.HandleFunc("/", server.uidHandler)
	return &server
}

func (server *Http) uidHandler(w http.ResponseWriter, r *http.Request) {
	hasUid := r.URL.Query().Has("uid")
	orderUid := r.URL.Query().Get("uid")
	if !hasUid {
		fmt.Fprintf(w, "uid value not specified\n")
		return
	}
	if order, ok := server.orderCash.Load(orderUid); ok {
		fmt.Fprintln(w, order)
	} else {
		fmt.Fprintf(w, "order with id=%s not found\n", orderUid)
	}
}

func (server *Http) Start() error {
	return server.ListenAndServe()
}
