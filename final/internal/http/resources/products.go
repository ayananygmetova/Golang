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

type ProductsResource struct {
	store  store.Store
	broker message_broker.MessageBroker
	cache  *lru.TwoQueueCache
}

func NewProductsResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *ProductsResource {
	return &ProductsResource{
		store:  store,
		broker: broker,
		cache:  cache,
	}
}

func (cr *ProductsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", cr.CreateProduct)
	r.Get("/", cr.AllProducts)
	r.Get("/{id}", cr.ById)
	r.Put("/", cr.UpdateProduct)
	r.Delete("/{id}", cr.DeleteProduct)

	return r
}

func (cr *ProductsResource) CreateProduct(w http.ResponseWriter, r *http.Request) {
	product := new(models.Product)
	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Products().Create(r.Context(), product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "BD err: %v", err)
		return
	}

	cr.broker.Cache().Purge() // в рамках учебного проекта полностью чистим кэш после создания новой категории

	w.WriteHeader(http.StatusCreated)
}

func (cr *ProductsResource) AllProducts(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.ProductsFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		productsFromCache, ok := cr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, productsFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	products, err := cr.store.Products().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" {
		cr.cache.Add(searchQuery, products)
	}

	render.JSON(w, r, products)
}

func (cr *ProductsResource) ById(w http.ResponseWriter, r *http.Request) {
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

	product, err := cr.store.Products().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.cache.Add(id, product)
	render.JSON(w, r, product)

}

func (cr *ProductsResource) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	product := new(models.Product)
	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(product,
		validation.Field(&product.ID, validation.Required),
		validation.Field(&product.Name, validation.Required))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err = cr.store.Products().Update(r.Context(), product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.broker.Cache().Remove(product.ID)
}

func (cr *ProductsResource) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err = cr.store.Products().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	cr.broker.Cache().Remove(id)
}
