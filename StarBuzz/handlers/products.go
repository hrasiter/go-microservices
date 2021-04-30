//	Package classification of Product API
//
//	Documentation for Product API
//
//	Schemes: http
//	Basepath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//	swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/microservices/StarBuzz/data"
)

//	A list of products returns in the response
//	swagger:response productResponse
type productResponse struct {
	//	All products in the system
	//	in: body
	Body []data.Product
}

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

type KeyProduct struct{}

//	swagger:route GET /products products listProducts
//	Returns a list of products
//	responses:
//		200:productsResponse
func (p Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling GET method")
	productList := data.GetProducts()
	err := productList.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to find Products", http.StatusInternalServerError)
		return
	}
}

func (p Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		p.l.Println("Cannot convert product id")
		http.Error(rw, "Product Not Found!", http.StatusNotFound)
		return
	}
	p.l.Println("Handling PUT method with id: ", id)

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err = data.UpdateProduct(id, &prod)

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

func (p Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handling POST method")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Printf("Product %#v", prod)
	data.AddProduct(&prod)

	prod.ToJSON(rw)
}

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	p.l.Println("Handling Middleware")
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("Deserializing product: ", err)
			http.Error(rw, "Error: reading product!", http.StatusBadRequest)
			return
		}

		err = prod.Validate()

		if err != nil {
			p.l.Println("Error validating product: ", err)

			http.Error(
				rw,
				fmt.Sprintf("Error: validating product: %s", err),
				http.StatusBadRequest,
			)

			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
