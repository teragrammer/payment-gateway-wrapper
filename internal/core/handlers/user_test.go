package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teragrammer/payment-gateway-wrapper/internal/config"
	"github.com/teragrammer/payment-gateway-wrapper/internal/core/middleware"
	"github.com/teragrammer/payment-gateway-wrapper/internal/core/validations"
	"github.com/teragrammer/payment-gateway-wrapper/internal/database/mongo"
	"github.com/teragrammer/payment-gateway-wrapper/internal/models"
	"github.com/teragrammer/payment-gateway-wrapper/internal/test/testutils"
	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
)

var userTest = struct {
	ID string
}{
	ID: "",
}

func TestCreateUserHandler(t *testing.T) {
	config.Load("../../../.env")
	db, _ := mongo.DefaultMongo()

	username, err := utils.GenerateRandomString(8)
	if err != nil {
		t.Fatalf("unable to generate random string for user password: %v", err)
	}
	inputs := map[string]interface{}{
		"full_name": "Doe Juan",
		"username":  username,
		"password":  "12345678",
		"role":      models.Manager,
	}
	query := map[string]string{}
	headers := map[string]string{
		"Authorization": "Bearer " + testutils.MockJWT("admin"),
	}

	// Set up the router
	url := "/api/v1/users"
	r := utils.SetupRouterPost(url, NewUserHandler(db).Create,
		middleware.JWT(), validations.ValidateUserInputs)

	rr, err := utils.SendPostRequest(r, url, inputs, query, headers)
	if err != nil {
		t.Fatalf("could not send post request to create a new user handler: %v", err)
	}

	assert.Equal(t, http.StatusOK, rr.Code)

	// Unmarshal the response body
	responseBody, err := utils.UnmarshalResponseBody(rr)
	if err != nil {
		t.Fatalf("could not unmarshal response from user handler: %v", err)
	}
	userTest.ID = responseBody["data"].(string)
}

func TestBrowseUsersHandler(t *testing.T) {
	db, _ := mongo.DefaultMongo()

	query := map[string]string{
		"search": "test",
	}
	headers := map[string]string{
		"Authorization": "Bearer " + testutils.MockJWT("admin"),
	}

	// Set up the router
	url := "/api/v1/users"
	r := utils.SetupRouterGet(url, NewUserHandler(db).Browse,
		middleware.JWT())

	rr, err := utils.SendGetRequest(r, url, query, headers)
	if err != nil {
		t.Fatalf("could not send get request to user handler: %v", err)
	}

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetUserHandler(t *testing.T) {
	db, _ := mongo.DefaultMongo()

	requestData := map[string]string{}
	headers := map[string]string{
		"Authorization": "Bearer " + testutils.MockJWT("admin"),
	}

	// Set up the router
	r := utils.SetupRouterGet("/api/v1/users/{id}", NewUserHandler(db).Get,
		middleware.JWT())

	rr, err := utils.SendGetRequest(r, "/api/v1/users/"+userTest.ID, requestData, headers)
	if err != nil {
		t.Fatalf("could not send get request to user handler: %v", err)
	}

	assert.Equal(t, http.StatusOK, rr.Code)

	responseBody, err := utils.UnmarshalResponseBody(rr)
	if err != nil {
		t.Fatalf("could not unmarshal response from user handler: %v", err)
	}
	assert.Equal(t, userTest.ID, responseBody["data"].(map[string]interface{})["_id"].(string))
}

func TestUpdateUserHandler(t *testing.T) {
	db, _ := mongo.DefaultMongo()

	username, err := utils.GenerateRandomString(8)
	inputs := map[string]interface{}{
		"full_name": "Joe Juan",
		"username":  username,
		"password":  "12345678",
		"role":      models.Manager,
	}
	query := map[string]string{}
	headers := map[string]string{
		"Authorization": "Bearer " + testutils.MockJWT("admin"),
	}

	// Set up the router
	r := utils.SetupRouterPut("/api/v1/users/{id}", NewUserHandler(db).Update,
		middleware.JWT(), validations.ValidateUserInputs)

	rr, err := utils.SendPutRequest(r, "/api/v1/users/"+userTest.ID, inputs, query, headers)
	if err != nil {
		t.Fatalf("could not send put request to user handler: %v", err)
	}

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteUserHandler(t *testing.T) {
	db, _ := mongo.DefaultMongo()

	query := map[string]string{}
	headers := map[string]string{
		"Authorization": "Bearer " + testutils.MockJWT("admin"),
	}

	// Set up the router
	r := utils.SetupRouterDelete("/api/v1/users/{id}", NewUserHandler(db).Delete,
		middleware.JWT())

	rr, err := utils.SendDeleteRequest(r, "/api/v1/users/"+userTest.ID, query, headers)
	if err != nil {
		t.Fatalf("could not send delete request to user handler: %v", err)
	}

	assert.Equal(t, http.StatusOK, rr.Code)
}
