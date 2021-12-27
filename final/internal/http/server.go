package http

import (
	"context"
	"final/internal/store"
	"log"
	"net/http"
	"time"

	lru "github.com/hashicorp/golang-lru"

	"final/internal/http/resources"
	"final/internal/message_broker"

	"github.com/go-chi/chi"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store

	cache   *lru.TwoQueueCache
	broker  message_broker.MessageBroker
	Address string
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	srv := &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()
	// r.Post("/products", func(w http.ResponseWriter, r *http.Request) {
	// 	product := new(models.Product)
	// 	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
	// 		w.WriteHeader(http.StatusUnprocessableEntity)
	// 		fmt.Fprintf(w, "Unknown err: %v", err)
	// 		return
	// 	}

	// 	if _, c_err := s.store.Categories().ByID(r.Context(), product.CategoryId); c_err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		fmt.Fprintf(w, "Category with id %d doesn't exist", product.CategoryId)
	// 		return
	// 	}

	// 	if err := s.store.Products().Create(r.Context(), product); err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		fmt.Fprintf(w, "DB err: %v", err)
	// 		return
	// 	}

	// 	w.WriteHeader(http.StatusCreated)
	// })
	// r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
	// 	products, err := s.store.Products().All(r.Context())
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Unknown err: %v", err)
	// 		return
	// 	}

	// 	render.JSON(w, r, products)
	// })
	// r.Get("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	idStr := chi.URLParam(r, "id")
	// 	id, err := strconv.Atoi(idStr)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		fmt.Fprintf(w, "Unknown err: %v", err)
	// 		return
	// 	}

	// 	product, err := s.store.Products().ByID(r.Context(), id)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		fmt.Fprintf(w, "Unknown err: %v", err)
	// 		return
	// 	}

	// 	render.JSON(w, r, product)
	// })
	// r.Put("/products", func(w http.ResponseWriter, r *http.Request) {
	// 	product := new(models.Product)
	// 	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
	// 		w.WriteHeader(http.StatusUnprocessableEntity)
	// 		fmt.Fprintf(w, "Unknown err: %v", err)
	// 		return
	// 	}
	// 	if _, c_err := s.store.Categories().ByID(r.Context(), product.CategoryId); c_err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		fmt.Fprintf(w, "Category with id %d doesn't exist", product.CategoryId)
	// 		return
	// 	}
	// 	err := validation.ValidateStruct(
	// 		product,
	// 		validation.Field(&product.ID, validation.Required),
	// 		validation.Field(&product.Name, validation.Required),
	// 		validation.Field(&product.CategoryId, validation.Required),
	// 	)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusUnprocessableEntity)
	// 		fmt.Fprintf(w, "Unknown err: %v", err)
	// 		return
	// 	}

	// 	if err := s.store.Products().Update(r.Context(), product); err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		fmt.Fprintf(w, "DB err: %v", err)
	// 		return
	// 	}
	// 	render.JSON(w, r, product)
	// })
	// r.Delete("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	idStr := chi.URLParam(r, "id")
	// 	id, err := strconv.Atoi(idStr)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		fmt.Fprintf(w, "Unknown err: %v", err)
	// 		return
	// 	}

	// 	if err := s.store.Products().Delete(r.Context(), id); err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		fmt.Fprintf(w, "DB err: %v", err)
	// 		return
	// 	}
	// })
	// r.Get("/categories/{id}/products", func(w http.ResponseWriter, r *http.Request) {
	// 	idStr := chi.URLParam(r, "id")
	// 	id, err := strconv.Atoi(idStr)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Unknown err: %v", err)
	// 		return
	// 	}
	// 	products, err := s.store.Products().ByCategory(r.Context(), id)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Unknown err: %v", err)
	// 		return
	// 	}

	// 	render.JSON(w, r, products)
	// })
	categoriesResource := resources.NewCategoriesResource(s.store, s.broker, s.cache)
	r.Mount("/categories", categoriesResource.Routes())

	productsResource := resources.NewProductsResource(s.store, s.broker, s.cache)
	r.Mount("/products", productsResource.Routes())
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
