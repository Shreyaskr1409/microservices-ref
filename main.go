package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	// there is a greedy-matching while matching the routes
	// this means any route which is not mentioned in any of
	// the HandleFunc will get routed to the function mentioned below
	// http.ResponseWriter implements io.Writer
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// r.Body implements io.ReadCloser
		d, err := io.ReadAll(r.Body)
		if err != nil {
			// w.WriteHeader(http.StatusBadRequest)
			// w.Write([]byte("Error in reading the body"))
			http.Error(w, "Error in reading the body", http.StatusBadRequest)
			return
		}
		// %s will turn bytes into string
		log.Printf("Hello %s\n", d)
		// Fprintf formats data into a io.Writer
		fmt.Fprintf(w, "Hello %s\n", d)
	})
	// the routes /bye/bye, /bye/world, etc. will not run the
	// function below but will run if we change the string addr from
	// "/bye" to "/bye/".
	// this is usually how greedy matching in ServeMux runs
	http.HandleFunc("/bye", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Goodbye World")
	})

	// ListenAndServe takes in a TCP network address and then calls
	// Serve with handler to handle requests on incoming connections
	// Unless someone specifier a new ServeMux, DefaultServeMux is
	// used to Serve (Serve here is a function in http).
	http.ListenAndServe(":9090", nil)
	// also remember that DefaultServeMux implements Handler interface
}
