package validations

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/teragrammer/payment-gateway-wrapper/internal/database/mongo"
	"github.com/teragrammer/payment-gateway-wrapper/internal/repository"
	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
)

func ValidateAuthenticationInputs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Authentication struct {
			Username string `json:"username" validate:"required,min=3,max=32"`
			Password string `json:"password" validate:"required,min=3,max=32"`
		}

		var form Authentication
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			utils.JSONErrorMessage(w, http.StatusBadRequest, "JSON", "Invalid JSON body")
			return
		}

		if !utils.Validate(w, form) {
			return
		}

		db, err := mongo.DefaultMongo()
		if err != nil {
			utils.JSONErrorMessage(w, http.StatusInternalServerError, "DB", "Unable to connect to database")
			return
		}

		user, err := repository.NewUserRepository(db).GetByUsername(form.Username)
		if err != nil {
			utils.JSONErrorValidation(w, map[string]string{
				"username": "Username not found",
			})
			return
		}

		if !utils.ValidatePassword(user.Password, form.Password) {
			utils.JSONErrorValidation(w, map[string]string{
				"password": "Invalid username or password",
			})
			return
		}

		// Attach the validated input to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
