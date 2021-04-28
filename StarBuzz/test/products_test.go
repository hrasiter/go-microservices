package test

import (
	"testing"

	"github.com/microservices/StarBuzz/data"
)

func TestCheckValidation(t *testing.T) {
	p := &data.Product{
		Name:  "Coffee",
		Price: 2.49,
		SKU:   "aaaaa-a-d",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
