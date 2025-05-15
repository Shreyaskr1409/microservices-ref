package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/Shreyaskr/microservices-ref/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	// encoding/json will be used to marshall a struct into json
	lp := data.GetProducts()

	// // We can either use Marshal or we can use Encoder
	// data, err := json.Marshal(lp)
	// if err != nil {
	// 	p.l.Println("Unable to marshal json")
	// 	http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	// }
	// w.Write(data)

	w.Header().Set("Content-Type", "application/json")

	// Encoder is faster than Marshal
	// Marshal send entire data at once while Encoder streams data
	if err := lp.ToJSON(w); err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) PostProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	data.AddProduct(prod)
}

func (p Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// we removed id from the parameters
	p.l.Println("Handle PUT Products")

	// mux.Vars() will contain the id which we extracted from the router
	// gorilla puts the extracted part inside the Request from where we
	// can extract those values
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}
}

// key to store in the values passed on through context of the request
type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		if err := prod.FromJSON(r.Body); err != nil {
			http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		// now we will put the product into Conext of the request
		// as every request has a context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
