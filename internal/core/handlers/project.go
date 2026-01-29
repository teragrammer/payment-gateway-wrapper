package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/teragrammer/payment-gateway-wrapper/internal"
	middleware2 "github.com/teragrammer/payment-gateway-wrapper/internal/core/middleware"
	"github.com/teragrammer/payment-gateway-wrapper/internal/core/validations"
	"github.com/teragrammer/payment-gateway-wrapper/internal/repository"
	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type ProjectHandler struct {
	bootstrap internal.Bootstrap
}

func NewProjectHandler(bootstrap internal.Bootstrap) *ProjectHandler {
	return &ProjectHandler{bootstrap: bootstrap}
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	inputs, ok := r.Context().Value("inputs").(bson.D)
	if !ok {
		utils.JSONErrorMessage(w, http.StatusInternalServerError, "CTX_INPUTS", "Error retrieving insert input")
		return
	}

	id, err := repository.NewProjectRepository(h.bootstrap.DB).Create(inputs)
	if err != nil {
		utils.JSONErrorMessage(w, http.StatusInternalServerError, "DB", "Unable to insert project")
		return
	}

	utils.JSONSuccess(w, http.StatusOK, id)
}

func (h *ProjectHandler) Browse(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	pageStr, limitStr := utils.ExtractPaginationFromRequest(r)

	result, err := repository.NewProjectRepository(h.bootstrap.DB).Browse(search, pageStr, limitStr)
	if err != nil {
		utils.JSONErrorMessage(w, http.StatusInternalServerError, "DB_ERROR", "Whoops something went wrong")
		return
	}

	utils.JSONSuccess(w, http.StatusOK, result)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	result, err := repository.NewProjectRepository(h.bootstrap.DB).GetByID(id)
	if err != nil {
		utils.JSONErrorMessage(w, http.StatusNotFound, "NOT_FOUND", "Unable to find project")
		log.Println(err.Error())
		return
	}

	utils.JSONSuccess(w, http.StatusOK, result)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	inputs, ok := r.Context().Value("inputs").(bson.D)
	if !ok {
		utils.JSONErrorMessage(w, http.StatusInternalServerError, "CTX_INPUTS", "Error retrieving update input")
		return
	}

	id := chi.URLParam(r, "id")
	err := repository.NewProjectRepository(h.bootstrap.DB).Update(id, inputs)
	if err != nil {
		utils.JSONErrorMessage(w, http.StatusInternalServerError, "DB", "Unable to update project")
		return
	}

	utils.WriteEmpty(w)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := repository.NewProjectRepository(h.bootstrap.DB).Delete(id)
	if err != nil {
		utils.JSONErrorMessage(w, http.StatusInternalServerError, "DB", "Unable to delete project")
		return
	}

	utils.WriteEmpty(w)
}

func ProjectRoutes(bootstrap internal.Bootstrap) func(r chi.Router) {
	handler := &ProjectHandler{bootstrap: bootstrap}
	return func(r chi.Router) {
		r.Use(middleware2.RateLimitIP())
		r.Use(middleware2.JWT())

		r.With(validations.ValidateProjectInputs).Post("/", handler.Create)
		r.Get("/", handler.Browse)
		r.Get("/{id}", handler.Get)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	}
}
