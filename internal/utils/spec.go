package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/go-chi/chi/v5"
)

func SetupRouterGet(route string, fn http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()
	for _, mw := range middlewares {
		r.Use(mw)
	}
	r.Get(route, fn)
	return r
}

func SendGetRequest(r *chi.Mux, url string, query map[string]string, headers map[string]string) (*httptest.ResponseRecorder, error) {
	return SendRequest(r, "GET", url, nil, query, headers)
}

func SetupRouterPost(route string, fn http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()
	for _, mw := range middlewares {
		r.Use(mw)
	}
	r.Post(route, fn)
	return r
}

func SendPostRequest(r *chi.Mux, url string, inputs interface{}, query map[string]string, headers map[string]string) (*httptest.ResponseRecorder, error) {
	return SendRequest(r, "POST", url, inputs, query, headers)
}

func SetupRouterPut(route string, fn http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()
	for _, mw := range middlewares {
		r.Use(mw)
	}
	r.Put(route, fn)
	return r
}

func SendPutRequest(r *chi.Mux, url string, inputs interface{}, query map[string]string, headers map[string]string) (*httptest.ResponseRecorder, error) {
	return SendRequest(r, "PUT", url, inputs, query, headers)
}

func SetupRouterDelete(route string, fn http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()
	for _, mw := range middlewares {
		r.Use(mw)
	}
	r.Delete(route, fn)
	return r
}

func SendDeleteRequest(r *chi.Mux, url string, query map[string]string, headers map[string]string) (*httptest.ResponseRecorder, error) {
	return SendRequest(r, "DELETE", url, nil, query, headers)
}

func SendRequest(r *chi.Mux, method string, url string, inputs interface{}, query map[string]string, headers map[string]string) (*httptest.ResponseRecorder, error) {
	methodSelector := http.MethodGet
	switch method {
	case "POST":
		methodSelector = http.MethodPost
	case "PUT":
		methodSelector = http.MethodPut
	case "DELETE":
		methodSelector = http.MethodDelete
	default:
		methodSelector = http.MethodGet
	}

	if len(query) > 0 {
		url = url + "?" + BuildQueryString(query)
	}

	var jsonData []byte
	var err error

	if inputs != nil {
		// Marshal request data to JSON
		jsonData, err = json.Marshal(inputs)
		if err != nil {
			return nil, err
		}
	}

	// Create a METHOD request with the JSON data
	req, err := http.NewRequest(methodSelector, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set Content-Type to application/json
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Record the HTTP response
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	return rr, nil
}

func UnmarshalResponseBody(rr *httptest.ResponseRecorder) (map[string]interface{}, error) {
	var responseBody map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
	return responseBody, err
}

func BuildQueryString(params map[string]string) string {
	// Create a URL-encoded query string from the map
	query := url.Values{}
	for key, value := range params {
		query.Set(key, value)
	}
	return query.Encode()
}
