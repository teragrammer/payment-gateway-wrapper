package handlers

import (
	"net/http"

	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
)

func Health(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status string `json:"status"`
	}

	resp := Response{
		Status: "ok",
	}

	utils.JSONSuccess(w, http.StatusOK, resp)
}
