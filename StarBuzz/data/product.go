package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {

	re := regexp.MustCompile("[a-z]+-[a-z]+-[a-z]+")

	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)

	if err != nil {
		return err
	}

	p.ID = id
	productList[pos] = p

	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
}

func getNextID() int {
	p := productList[len(productList)-1]
	return p.ID + 1
}

func GetProducts() Products {
	return productList
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "aha453",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee",
		Price:       1.99,
		SKU:         "abc123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
