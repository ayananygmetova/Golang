package resources

import (
	"encoding/json"
	"fmt"
	"hw9/internal/cache"
	"hw9/internal/models"
	"hw9/internal/store"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
)

type ProductsResource struct {
	store store.Store
	cache cache.Cache
}

func NewProductsResource(store store.Store, cache cache.Cache) *ProductsResource {
	return &ProductsResource{
		store: store,
		cache: cache,
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

	if err := cr.cache.DeleteAll(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Cache err: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cr *ProductsResource) AllProducts(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.ProductsFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		productsFromCache, err := cr.cache.Products().Get(r.Context(), searchQuery)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Cache err: %v", err)
			return
		}
		if productsFromCache != nil {
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

	if searchQuery != "" && len(products) > 0 {
		err = cr.cache.Products().Set(r.Context(), searchQuery, products)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Cache err: %v", err)
			return
		}
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

	product, err := cr.store.Products().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

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
}
