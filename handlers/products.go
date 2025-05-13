package handlers

import (
	"log"
	"net/http"

	"github.com/Shreyaskr/microservices-ref/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.getProducts(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	// encoding/json will be used to marshall a struct into json
	lp := data.GetProducts()

	// // We can either use Marshal or we can use Encoder
	// data, err := json.Marshal(lp)
	// if err != nil {
	// 	p.l.Println("Unable to marshal json")
	// 	http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	// }
	// w.Write(data)

	// Encoder is faster than Marshal
	// Marshal send entire data at once while Encoder streams data
	if err := lp.ToJSON(w); err != nil {
		p.l.Println("Unable to marshal json")
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}
