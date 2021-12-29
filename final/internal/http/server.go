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
	categoriesResource := resources.NewCategoriesResource(s.store, s.broker, s.cache)
	r.Mount("/categories", categoriesResource.Routes())

	productsResource := resources.NewProductsResource(s.store, s.broker, s.cache)
	r.Mount("/products", productsResource.Routes())

	propertiesResource := resources.NewPropertyResource(s.store, s.broker, s.cache)
	r.Mount("/properties", propertiesResource.Routes())

	characteristicsResource := resources.NewCharacteristiccResource(s.store, s.broker, s.cache)
	r.Mount("/characteristics", characteristicsResource.Routes())

	productCharacteristicsResource := resources.NewProductCharacteristicsResource(s.store, s.broker, s.cache)
	r.Mount("/products/{id}/characteristics", productCharacteristicsResource.Routes())

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
