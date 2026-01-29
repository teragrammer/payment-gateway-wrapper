package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teragrammer/payment-gateway-wrapper/internal/config"
	"github.com/teragrammer/payment-gateway-wrapper/internal/core/validations"
	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
)

func TestLoginHandler(t *testing.T) {
	config.Load("../../../.env")

	inputs := map[string]interface{}{
		"username": "admin",
		"password": "12345678",
	}
	query := map[string]string{}
	headers := map[string]string{}

	// Set up the router
	url := "/api/v1/login"
	r := utils.SetupRouterPost(url, Login,
		validations.ValidateAuthenticationInputs)

	rr, err := utils.SendPostRequest(r, url, inputs, query, headers)
	if err != nil {
		t.Fatalf("could not send post request: %v", err)
	}

	assert.Equal(t, http.StatusOK, rr.Code)
}
