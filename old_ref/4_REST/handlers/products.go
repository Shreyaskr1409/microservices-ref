package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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

	case http.MethodPost:
		p.postProduct(w, r)

	case http.MethodPut:
		p.l.Println("PUT: ", r.URL.Path)
		// we will put data into an ID, so expect ID in the URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URL, more than one id")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			p.l.Println("Invalid URL, more than one capture group")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		// since we have taken the string from regex, the string
		// will surely have a number in it
		id, _ := strconv.Atoi(idString)

		p.updateProduct(id, w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
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

	// Encoder is faster than Marshal
	// Marshal send entire data at once while Encoder streams data
	if err := lp.ToJSON(w); err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) postProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")

	prod := &data.Product{}
	if err := prod.FromJSON(r.Body); err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Products")

	prod := &data.Product{}
	if err := prod.FromJSON(r.Body); err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err := data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}
}
