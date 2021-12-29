package resources

import (
	"encoding/json"
	"final/internal/message_broker"
	"final/internal/models"
	"final/internal/store"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	lru "github.com/hashicorp/golang-lru"
)

type ProductCharacteristicsResource struct {
	store  store.Store
	broker message_broker.MessageBroker
	cache  *lru.TwoQueueCache
}

func NewProductCharacteristicsResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *ProductCharacteristicsResource {
	return &ProductCharacteristicsResource{
		store:  store,
		broker: broker,
		cache:  cache,
	}
}

func (cr *ProductCharacteristicsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/characteristics", cr.CreateProductCharacteristics)
	r.Get("/{id}/characteristics", cr.ById)
	r.Delete("/{id}/characteristics/{characteristics_id}", cr.DeleteProductCharacteristics)

	return r
}

func (cr *ProductCharacteristicsResource) CreateProductCharacteristics(w http.ResponseWriter, r *http.Request) {
	p_c := new(models.ProductCharacteristics)
	if err := json.NewDecoder(r.Body).Decode(p_c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.ProductCharacteristics().Create(r.Context(), p_c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "BD err: %v", err)
		return
	}

	cr.broker.Cache().Purge() // в рамках учебного проекта полностью чистим кэш после создания новой категории

	w.WriteHeader(http.StatusCreated)
}

func (cr *ProductCharacteristicsResource) ById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	productFromCache, ok := cr.cache.Get(id)
	if ok {
		render.JSON(w, r, productFromCache)
		return
	}

	products, err := cr.store.ProductCharacteristics().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	render.JSON(w, r, products)

}

func (cr *ProductCharacteristicsResource) DeleteProductCharacteristics(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err = cr.store.ProductCharacteristics().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.broker.Cache().Remove(id)
}
