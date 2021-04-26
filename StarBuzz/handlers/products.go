package handlers

import (
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

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.GetProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling GET method")
	productList := data.GetProducts()
	err := productList.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to find Products", http.StatusInternalServerError)
		return
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling POST method")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to store new product!", http.StatusBadRequest)
	}

	p.l.Printf("Product %#v", prod)
	data.AddProduct(prod)

	prod.ToJSON(rw)
}
