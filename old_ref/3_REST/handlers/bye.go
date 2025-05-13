package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type Bye struct {
	l *log.Logger
}

func NewBye(l *log.Logger) *Bye {
	return &Bye{l}
}

func (b *Bye) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.l.Println("Bye World")
	fmt.Fprintln(w, "Byebyee")
}
