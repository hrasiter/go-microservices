package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/microservices/StarBuzz/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (*Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	productList := data.GetProducts()
	d, err := json.Marshal(productList)

	if err != nil {
		http.Error(rw, "Unable to find Products", http.StatusInternalServerError)
		return
	}

	rw.Write(d)

}
