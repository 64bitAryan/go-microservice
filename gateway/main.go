package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/64bitAryan/go-microservice/aggregator/client"
	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type InvoiceHandler struct {
	client client.Client
}

func newInvoiceHandler(client client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: client,
	}
}

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "listen adderess of the gateway")
	aggregatorServiceAddr := flag.String("aggServiceAddr", "http://localhost:3000", "listen adderess of aggregator service")
	flag.Parse()
	var (
		client     = client.NewHTTPClient(*aggregatorServiceAddr)
		invHandler = newInvoiceHandler(client)
	)
	http.HandleFunc("/invoice", makeApiFunc(invHandler.handleGetInvoice))
	logrus.Printf("gateway is runing on port %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	inv, err := h.client.GetInvoice(context.Background(), 4931244646797028242)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeApiFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
