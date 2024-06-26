package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rest_api_postgres_clean/internal/inventory"
	"strings"
)

func (h *HTTPserver) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/products/"):]
	if id == "" || strings.ContainsRune(id, '/') {
		http.NotFound(w, r)
		return

	}
	product, err := h.inventory.GetProducts(r.Context(), id)
	switch {
	case err == context.Canceled, err == context.DeadlineExceeded:
		{
			return
		}
	case err != nil:
		{
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Println("error while getting product :::", err)
		}
	case product == nil:
		http.Error(w, "Product ot found", http.StatusNotFound)

	default:
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "\t")
		err := enc.Encode(product)
		if err != nil {
			log.Println("Error while encoding product ::: ", err)
		}
	}
}

func NewHttpServer(service *inventory.Service) http.Handler {
	s := HTTPserver{
		inventory: service,
		mux:       http.NewServeMux(),
	}

	s.mux.HandleFunc("/product/", s.handleGetProduct)
	return s.mux
}

// Server expose inventory
type HTTPserver struct {
	inventory *inventory.Service
	mux       *http.ServeMux
}

func (s *httpServer) Shutdown(ctx context.Context) {
	log.Println("shutting down HTTP server")
	if s.server != nil {
		if err := s.server.Shutdown(ctx); err != nil {
			log.Println("graceful shutdown of HTTP server failed")
		}
	}
}
