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

type CharacteristicsResource struct {
	store  store.Store
	broker message_broker.MessageBroker
	cache  *lru.TwoQueueCache
}

func NewCharacteristiccResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *CharacteristicsResource {
	return &CharacteristicsResource{
		store:  store,
		broker: broker,
		cache:  cache,
	}
}

func (cr *CharacteristicsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", cr.CreateCharacteristics)
	r.Get("/", cr.AllCharacteristics)
	r.Get("/{id}", cr.ById)
	r.Put("/", cr.UpdateCharacteristics)
	r.Delete("/{id}", cr.DeleteCharacteristics)

	return r
}

func (cr *CharacteristicsResource) CreateCharacteristics(w http.ResponseWriter, r *http.Request) {
	characteristic := new(models.Characteristics)
	if err := json.NewDecoder(r.Body).Decode(characteristic); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Characteristics().Create(r.Context(), characteristic); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "BD err: %v", err)
		return
	}

	cr.broker.Cache().Purge() // в рамках учебного проекта полностью чистим кэш после создания новой категории

	w.WriteHeader(http.StatusCreated)
}

func (cr *CharacteristicsResource) AllCharacteristics(w http.ResponseWriter, r *http.Request) {
	characteristics, err := cr.store.Characteristics().All(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
	render.JSON(w, r, characteristics)
}

func (cr *CharacteristicsResource) ById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	characteristicFromCache, ok := cr.cache.Get(id)
	if ok {
		render.JSON(w, r, characteristicFromCache)
		return
	}

	characteristic, err := cr.store.Characteristics().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.cache.Add(id, characteristic)
	render.JSON(w, r, characteristic)

}

func (cr *CharacteristicsResource) UpdateCharacteristics(w http.ResponseWriter, r *http.Request) {
	characteristic := new(models.Characteristics)
	if err := json.NewDecoder(r.Body).Decode(characteristic); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(characteristic,
		validation.Field(&characteristic.ID, validation.Required),
		validation.Field(&characteristic.Value, validation.Required))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err = cr.store.Characteristics().Update(r.Context(), characteristic); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.broker.Cache().Remove(characteristic.ID)
}

func (cr *CharacteristicsResource) DeleteCharacteristics(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err = cr.store.Characteristics().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.broker.Cache().Remove(id)
}
