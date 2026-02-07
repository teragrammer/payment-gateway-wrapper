package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	middleware2 "github.com/teragrammer/payment-gateway-wrapper/internal/core/middleware"
	"github.com/teragrammer/payment-gateway-wrapper/internal/core/validations"
	"github.com/teragrammer/payment-gateway-wrapper/internal/repository"
	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	db *mongo.Database
}

func NewUserHandler(db *mongo.Database) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	inputs, ok := r.Context().Value("inputs").(bson.D)
	if !ok {
		utils.JSONErrorMessage(w, http.StatusInternalServerError, "CTX_INPUTS", "Error retrieving insert input")
		return
	}

	id, err := repository.NewUserRepository(h.db).Create(inputs)
	if err != nil {
		utils.JSONErrorMessage(w, http.StatusInternalServerError, "DB", "Unable to insert user")
		return
	}

	utils.JSONSuccess(w, http.StatusOK, id)
}

func (h *UserHandler) Browse(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	pageStr, limitStr := utils.ExtractPaginationFromRequest(r)

	result, err := repository.NewUserRepository(h.db).Browse(search, pageStr, limitStr)
	if err != nil {
		utils.JSONErrorMessage(w, http.StatusInternalServerError, "DB_ERROR", "Whoops something went wrong")
		return
	}

	utils.JSONSuccess(w, http.StatusOK, result)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	result, err := repository.NewUserRepository(h.db).GetByID(id)
	if err != nil {
		utils.JSONErrorMessage(w, http.StatusNotFound, "NOT_FOUND", "Unable to find user")
		log.Println(err.Error())
		return
	}

	utils.JSONSuccess(w, http.StatusOK, result)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	inputs, ok := r.Context().Value("inputs").(bson.D)
	if !ok {
		utils.JSONErrorMessage(w, http.StatusInternalServerError, "CTX_INPUTS", "Error retrieving update input")
		return
	}

	id := chi.URLParam(r, "id")
	err := repository.NewUserRepository(h.db).Update(id, inputs)
	if err != nil {
		utils.JSONErrorMessage(w, http.StatusInternalServerError, "DB", "Unable to update user")
		return
	}

	utils.WriteEmpty(w)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := repository.NewUserRepository(h.db).Delete(id)
	if err != nil {
		utils.JSONErrorMessage(w, http.StatusInternalServerError, "DB", "Unable to delete user")
		return
	}

	utils.WriteEmpty(w)
}

func UserRoutes(db *mongo.Database) func(r chi.Router) {
	handler := NewUserHandler(db)
	return func(r chi.Router) {
		r.Use(middleware2.RateLimitIP())
		r.Use(middleware2.JWT())

		r.With(validations.ValidateUserInputs).Post("/", handler.Create)
		r.Get("/", handler.Browse)
		r.Get("/{id}", handler.Get)
		r.With(validations.ValidateUserInputs).Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	}
}
