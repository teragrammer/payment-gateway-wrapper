package validations

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/teragrammer/payment-gateway-wrapper/internal/services"
	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func ValidateProjectInputs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := services.ExtractClaims(r.Context())
		if !ok {
			utils.JSONErrorMessage(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired authentication token")
			return
		}

		type Project struct {
			Name     string `json:"name" validate:"required,min=3,max=100"`
			IsActive int    `json:"is_active" validate:"required,oneof=0 1"`
		}

		var form Project
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			utils.JSONErrorMessage(w, http.StatusBadRequest, "JSON", "Invalid JSON body")
			return
		}

		if !utils.Validate(w, form) {
			return
		}

		publicKey, err := utils.GenerateKey(32)
		if err != nil {
			utils.JSONErrorMessage(w, http.StatusBadRequest, "PUB_KEY", "Unable to generate public key")
			return
		}

		privateKey, err := utils.GenerateKey(64)
		if err != nil {
			utils.JSONErrorMessage(w, http.StatusBadRequest, "PRI_KEY", "Unable to generate private key")
			return
		}

		// cleaned inputs
		inputs := bson.D{
			{"name", form.Name},
			{"user_id", claims.Id},
			{"public_key", publicKey},
			{"private_key", privateKey},
			{"is_activate", form.IsActive},
		}

		// create
		if r.Method == http.MethodPost {
			inputs = append(inputs, bson.E{Key: "created_at", Value: time.Now()})
		}

		// update
		if r.Method == http.MethodPut {
			inputs = append(inputs, bson.E{Key: "updated_at", Value: time.Now()})
		}

		// attach the validated input to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "inputs", inputs)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
