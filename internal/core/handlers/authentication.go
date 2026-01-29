package handlers

import (
	"net/http"

	"github.com/teragrammer/payment-gateway-wrapper/internal/models"
	"github.com/teragrammer/payment-gateway-wrapper/internal/services"
	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		utils.JSONErrorMessage(w, http.StatusInternalServerError, "CTX_USER", "Error retrieving user")
		return
	}

	token, err := services.GenerateAccessToken(user.ID.Hex(), user.Role)
	if err != nil {
		utils.JSONErrorMessage(w, http.StatusInternalServerError, "TOKEN", "Unable to generate access token")
		return
	}

	utils.JSONSuccess(w, http.StatusOK, token)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	_, ok := services.ExtractClaims(r.Context())
	if !ok {
		utils.JSONErrorMessage(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired authentication token")
		return
	}

	utils.WriteEmpty(w)
}
