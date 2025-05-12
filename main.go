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
)

func main() {
	// Logger takes in io.Writer, prefix string and log flag
	l := log.New(os.Stdout, "product-api: ", log.LstdFlags)
	hh := handlers.NewHello(l)
	gh := handlers.NewBye(l)

	// Given below is the way Handler is implemented
	//
	// type Handler interface {
	//     ServeHTTP(ResponseWriter, *Request)
	// }
	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/bye", gh)

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
	go func() {
		err := s.ListenAndServe()
		if err != nil {
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
}
