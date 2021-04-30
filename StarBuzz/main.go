package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/microservices/StarBuzz/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	ph := handlers.NewProducts(l)

	router := mux.NewRouter()
	getRouter := router.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := router.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := router.Methods("POST").Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareValidateProduct)

	opts := middleware.RedocOpts{SpecURL: "./swagger.yaml"}
	sh := opts.Recdoc(opts, nil)

	getRouter.Handle("/docs", sh)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			l.Fatal(err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan

	l.Println("Received terminate, graceful shutdown", sig)

	tc, cf := context.WithTimeout(context.Background(), 30*time.Second)
	defer cf()
	s.Shutdown(tc)

}
