package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger // allows nice testing by assigning a logger specifically
}

// DEPENDENCY INJECTION
// Injects the logger into the Hello struct instead of using or creating
// it's own logger inside NewHello. This enables ease in testing where
// logger can be relaced with something similar
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// r.Body implements io.ReadCloser
	d, err := io.ReadAll(r.Body)
	if err != nil {
		// w.WriteHeader(http.StatusBadRequest)
		// w.Write([]byte("Error in reading the body"))
		http.Error(w, "Error in reading the body", http.StatusBadRequest)
		return
	}
	// %s will turn bytes into string
	h.l.Println("Hello World")
	// Fprintf formats data into a io.Writer
	fmt.Fprintf(w, "Hello %s\n", d)
}
