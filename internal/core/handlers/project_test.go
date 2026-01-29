package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teragrammer/payment-gateway-wrapper/internal"
	"github.com/teragrammer/payment-gateway-wrapper/internal/config"
	"github.com/teragrammer/payment-gateway-wrapper/internal/core/middleware"
	"github.com/teragrammer/payment-gateway-wrapper/internal/core/validations"
	"github.com/teragrammer/payment-gateway-wrapper/internal/mocks"
	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
)

var projectTest = struct {
	ID string
}{
	ID: "",
}

func TestCreateProjectHandler(t *testing.T) {
	config.Load("../../../.env")
	bootstrap := internal.InitializedBootstrap()

	// Set up the router
	url := "/api/v1/projects"
	r := utils.SetupRouterPost(url, NewProjectHandler(bootstrap).Create,
		middleware.JWT(), validations.ValidateProjectInputs)

	rr, err := utils.SendPostRequest(r, url, map[string]interface{}{
		"name":      "Test",
		"is_active": 1,
	}, map[string]string{}, map[string]string{
		"Authorization": "Bearer " + mocks.MockJWT("admin"),
	})
	if err != nil {
		t.Fatalf("could not send post request to project handler: %v", err)
	}

	assert.Equal(t, http.StatusOK, rr.Code)

	// Unmarshal the response body
	responseBody, err := utils.UnmarshalResponseBody(rr)
	if err != nil {
		t.Fatalf("could not unmarshal response from create project handler: %v", err)
	}
	projectTest.ID = responseBody["data"].(string)
}

func TestBrowseProjectsHandler(t *testing.T) {
	bootstrap := internal.InitializedBootstrap()

	query := map[string]string{
		"search": "Test",
	}
	headers := map[string]string{
		"Authorization": "Bearer " + mocks.MockJWT("admin"),
	}

	// Set up the router
	url := "/api/v1/projects"
	r := utils.SetupRouterGet(url, NewProjectHandler(bootstrap).Browse,
		middleware.JWT())

	rr, err := utils.SendGetRequest(r, url, query, headers)
	if err != nil {
		t.Fatalf("could not send get request to project handler: %v", err)
	}

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetProjectHandler(t *testing.T) {
	bootstrap := internal.InitializedBootstrap()

	requestData := map[string]string{}
	headers := map[string]string{
		"Authorization": "Bearer " + mocks.MockJWT("admin"),
	}

	// Set up the router
	r := utils.SetupRouterGet("/api/v1/projects/{id}", NewProjectHandler(bootstrap).Get,
		middleware.JWT())

	rr, err := utils.SendGetRequest(r, "/api/v1/projects/"+projectTest.ID, requestData, headers)
	if err != nil {
		t.Fatalf("could not send get request for project handler: %v", err)
	}

	assert.Equal(t, http.StatusOK, rr.Code)

	responseBody, err := utils.UnmarshalResponseBody(rr)
	if err != nil {
		t.Fatalf("could not unmarshal response from project response: %v", err)
	}
	assert.Equal(t, projectTest.ID, responseBody["data"].(map[string]interface{})["_id"].(string))
}

func TestUpdateProjectHandler(t *testing.T) {
	bootstrap := internal.InitializedBootstrap()

	inputs := map[string]interface{}{
		"name":      "New Test",
		"is_active": 1,
	}
	query := map[string]string{}
	headers := map[string]string{
		"Authorization": "Bearer " + mocks.MockJWT("admin"),
	}

	// Set up the router
	r := utils.SetupRouterPut("/api/v1/projects/{id}", NewProjectHandler(bootstrap).Update,
		middleware.JWT(), validations.ValidateProjectInputs)

	rr, err := utils.SendPutRequest(r, "/api/v1/projects/"+projectTest.ID, inputs, query, headers)
	if err != nil {
		t.Fatalf("could not send put request to project handler: %v", err)
	}

	fmt.Println(rr)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteProjectHandler(t *testing.T) {
	bootstrap := internal.InitializedBootstrap()

	query := map[string]string{}
	headers := map[string]string{
		"Authorization": "Bearer " + mocks.MockJWT("admin"),
	}

	// Set up the router
	r := utils.SetupRouterDelete("/api/v1/projects/{id}", NewProjectHandler(bootstrap).Delete,
		middleware.JWT())

	rr, err := utils.SendDeleteRequest(r, "/api/v1/projects/"+projectTest.ID, query, headers)
	if err != nil {
		t.Fatalf("could not send delete request to project handler: %v", err)
	}

	assert.Equal(t, http.StatusOK, rr.Code)
}
