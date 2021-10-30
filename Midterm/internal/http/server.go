package http

import (
	"Midterm/internal/models"
	"Midterm/internal/store"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store

	Address string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
		store:       store,

		Address: address,
	}
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	r.Post("/products", func(w http.ResponseWriter, r *http.Request) {
		product := new(models.Product)
		if err := json.NewDecoder(r.Body).Decode(product); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		_, err := s.store.Categories().ByID(r.Context(), product.CategoryId)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		s.store.Products().Create(r.Context(), product)
		render.JSON(w, r, product)
	})
	r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		products, err := s.store.Products().All(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, products)
	})
	r.Get("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		product, err := s.store.Products().ByID(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, product)
	})
	r.Put("/products", func(w http.ResponseWriter, r *http.Request) {
		product := new(models.Product)
		if err := json.NewDecoder(r.Body).Decode(product); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		_, err := s.store.Categories().ByID(r.Context(), product.CategoryId)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		s.store.Products().Update(r.Context(), product)
		render.JSON(w, r, product)
	})
	r.Delete("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Products().Delete(r.Context(), id)
	})
	r.Get("/categories/{id}/products", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		products, err := s.store.Products().ByCategory(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, products)
	})
	r.Post("/categories", func(w http.ResponseWriter, r *http.Request) {
		category := new(models.Category)
		if err := json.NewDecoder(r.Body).Decode(category); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Categories().Create(r.Context(), category)
		render.JSON(w, r, category)
	})
	r.Get("/categories", func(w http.ResponseWriter, r *http.Request) {
		categories, err := s.store.Categories().All(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, categories)
	})
	r.Get("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		category, err := s.store.Categories().ByID(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, category)
	})
	r.Put("/categories", func(w http.ResponseWriter, r *http.Request) {
		category := new(models.Category)
		if err := json.NewDecoder(r.Body).Decode(category); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Categories().Update(r.Context(), category)
		render.JSON(w, r, category)
	})
	r.Delete("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Categories().Delete(r.Context(), id)
	})

	return r
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(srv)

	log.Println("[HTTP] Server running on", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got err while shutting down^ %v", err)
	}

	log.Println("[HTTP] Proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	<-s.idleConnsCh
}
