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
	validation "github.com/go-ozzo/ozzo-validation"
	lru "github.com/hashicorp/golang-lru"
)

type PropertyResource struct {
	store  store.Store
	broker message_broker.MessageBroker
	cache  *lru.TwoQueueCache
}

func NewPropertyResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *PropertyResource {
	return &PropertyResource{
		store:  store,
		broker: broker,
		cache:  cache,
	}
}

func (cr *PropertyResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", cr.CreateProperty)
	r.Get("/", cr.AllProperties)
	r.Get("/{id}", cr.ById)
	r.Put("/", cr.UpdateProperty)
	r.Delete("/{id}", cr.DeleteProperty)

	return r
}

func (cr *PropertyResource) CreateProperty(w http.ResponseWriter, r *http.Request) {
	property := new(models.Property)
	if err := json.NewDecoder(r.Body).Decode(property); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Properties().Create(r.Context(), property); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "BD err: %v", err)
		return
	}

	cr.broker.Cache().Purge() // в рамках учебного проекта полностью чистим кэш после создания новой категории

	w.WriteHeader(http.StatusCreated)
}

func (cr *PropertyResource) AllProperties(w http.ResponseWriter, r *http.Request) {
	properties, err := cr.store.Properties().All(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	render.JSON(w, r, properties)
}

func (cr *PropertyResource) ById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	propertyFromCache, ok := cr.cache.Get(id)
	if ok {
		render.JSON(w, r, propertyFromCache)
		return
	}

	property, err := cr.store.Properties().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.cache.Add(id, property)
	render.JSON(w, r, property)

}

func (cr *PropertyResource) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	property := new(models.Property)
	if err := json.NewDecoder(r.Body).Decode(property); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(property,
		validation.Field(&property.ID, validation.Required),
		validation.Field(&property.Name, validation.Required))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err = cr.store.Properties().Update(r.Context(), property); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.broker.Cache().Remove(property.ID)
}

func (cr *PropertyResource) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err = cr.store.Properties().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.broker.Cache().Remove(id)
}
