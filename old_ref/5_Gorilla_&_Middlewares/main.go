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
	"github.com/gorilla/mux"
)

func main() {
	// Logger takes in io.Writer, prefix string and log flag
	l := log.New(os.Stdout, "product-api: ", log.LstdFlags)
	l.Println("Program started")

	ph := handlers.NewProducts(l)

	// we will replace the old code with gorilla framework
	// gorilla has the concepts of router and subrouters
	sm := mux.NewRouter()
	// from sm, the GET method is taken and converted into
	// a subrouter which can be used to configure stuffs
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)
	// sm.Handle("/products", ph)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.PostProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// ListenAndServe takes in a TCP network address and then calls
	// Serve with handler to handle requests on incoming connections
	// Unless someone specifier a new ServeMux, DefaultServeMux is
	// used to Serve (Serve here is a function in http).
	l.Println("Server started")
	go func() {
		err := s.ListenAndServe()
		// http.ErrServerClosed will be thrown when server is shutdown
		// will add to l.Fatal(err) if not handled
		// we will handle it separately instead at the end
		if err != nil && err != http.ErrServerClosed {
			l.Fatal(err)
		}
	}()
	// also remember that DefaultServeMux implements Handler interface

	// -------------- BUFFERED vs/ UNBUFFERED CHANNELS ---------------
	// Here we will get a warning for using an unbuffered channel.
	// The reason for this is that goroutines take a bit
	// of delay to check for channels. If the the goroutine is not
	// immediately ready to handle those signals, the signal will get
	// lost when we encounter another signal.
	// while this is important for other implementations for the chan,
	// we can move on with this since we are constantly blocking at few
	// lines below.
	// If we were using this for a 'for' loop, there might be sometimes
	// when some function is happening and some time when the channel
	// is checked
	// if we had to use buffered channel (which we stil can), we would
	// use:
	// sigChan := make(chan os.Signal, 1)
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	l.Println("Recieved terminate, graceful shutdown due to signal: ", sig)

	// timeout-context
	// allows 30s to close all handlers
	// tc is like a counter
	// cancel is used to release resources
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// shutdown everything forcefully if it is still working
	l.Println("Shutdown started")
	s.Shutdown(tc)
	l.Println("Shutdown successful")
}
