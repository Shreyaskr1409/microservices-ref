package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shreyaskr/microservices-ref/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "product-api: ", log.LstdFlags)
	l.Println("Program started")

	ph := handlers.NewProducts(l)

	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.PostProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)

    // GO-SWAGGER tooling acts completely broken as of right now
    // I would rather use other alternatives

    // // Redoc generates documentations using swagger o/p files
    // // RedocOpts will take file from url "/swagger.yaml" which is
    // // a route (not the file itself)
    // options := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
    // sh := middleware.Redoc(options, nil)
    // getRouter.Handle("/docs", sh)
    // // if we GET /docs, we will get 404 error regarding Redoc get a
    // // 404 error for the route /swagger.yaml
    // // to settle this, we need to serve our file: swagger.yaml at
    // // the route: /swagger.yaml
    // getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	l.Println("Server started")
	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			l.Fatal(err)
		}
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	l.Println("Recieved terminate, graceful shutdown due to signal: ", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	l.Println("Shutdown started")
	s.Shutdown(tc)
	l.Println("Shutdown successful")
}
