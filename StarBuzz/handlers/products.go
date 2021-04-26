package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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

	if r.Method == http.MethodPut {
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, error := strconv.Atoi(idString)

		if error != nil {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.l.Println("id received: ", id)
		p.updateProducts(id, rw, r)
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

func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling PUT method")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if err != nil {
		p.l.Println("Cannot create product from JSON")
		http.Error(rw, "Unable to find Products", http.StatusBadRequest)
		return
	}

	err = data.UpdateProduct(id, prod)

	if err == data.ErrProductNotFound {
		p.l.Println("Cannot find product")
		http.Error(rw, "Product Not Found!", http.StatusNotFound)
		return
	}

	if err != nil {
		p.l.Println("Internal Server Error")
		http.Error(rw, "Product Not Found!", http.StatusInternalServerError)
		return
	}

	prod.ToJSON(rw)

	p.l.Println("Product updated")

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
