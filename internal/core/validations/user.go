package validations

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/teragrammer/payment-gateway-wrapper/internal/database/mongo"
	"github.com/teragrammer/payment-gateway-wrapper/internal/repository"
	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func ValidateUserInputs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type User struct {
			FullName string `json:"full_name" validate:"required,min=2,max=100"`
			Username string `json:"username" validate:"required,min=6,max=32"`
			Password string `json:"password" validate:"required,min=8,max=32"`
			Role     int    `json:"role" validate:"required,oneof=0 1"`
		}

		var form User
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			utils.JSONErrorMessage(w, http.StatusBadRequest, "JSON", "Invalid JSON body")
			return
		}

		if !utils.Validate(w, form) {
			return
		}

		password, err := utils.GeneratePasswordHash(form.Password)
		if err != nil {
			utils.JSONErrorValidation(w, map[string]string{
				"password": "Unable to generate password hash",
			})
			return
		}

		// cleaned inputs
		inputs := bson.D{
			{"full_name", form.FullName},
			{"username", form.Username},
			{"password", password},
			{"role", form.Role},
		}

		db, err := mongo.DefaultMongo()
		if err != nil {
			utils.JSONErrorMessage(w, http.StatusInternalServerError, "DB", "Unable to connect to database")
			return
		}

		// create
		if r.Method == http.MethodPost {
			inputs = append(inputs, bson.E{Key: "created_at", Value: time.Now()})
		}

		// update
		if r.Method == http.MethodPut {
			id := chi.URLParam(r, "id")
			_, err := repository.NewUserRepository(db).GetByUniqueUsername(form.Username, id)
			if err == nil {
				utils.JSONErrorValidation(w, map[string]string{
					"username": "Username is already taken",
				})
				return
			}

			inputs = append(inputs, bson.E{Key: "updated_at", Value: time.Now()})
		}

		// attach the validated input to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "inputs", inputs)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
