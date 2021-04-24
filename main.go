package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Hello World")

		d, err := ioutil.ReadAll(r.Body)

		if err != nil {
			// rw.WriteHeader(http.StatusBadRequest)
			// rw.Write([]byte("Oooops"))
			http.Error(rw, "Ooops", http.StatusBadRequest)
			return
		}

		log.Printf("Data is %s\n", d)

		fmt.Fprintf(rw, "Hello your data is %s\n", d)
	})

	http.HandleFunc("/goodbye", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Goodbye World")
	})

	http.ListenAndServe(":9090", nil)
}
