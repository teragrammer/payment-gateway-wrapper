package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

type SuccessResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data,omitempty"`
}

type ErrorMessageResponse struct {
	Success   bool   `json:"success"`
	Code      string `json:"code"`
	Details   string `json:"details,omitempty"`
	Timestamp string `json:"timestamp"`
}

type ErrorObjectResponse struct {
	Success   bool              `json:"success"`
	Code      string            `json:"code"`
	Message   string            `json:"message,omitempty"`
	Errors    map[string]string `json:"errors,omitempty"`
	Timestamp string            `json:"timestamp"`
}

func JSONSuccess[T any](w http.ResponseWriter, status int, data T) {
	response := SuccessResponse[T]{
		Success: true,
		Data:    data,
	}

	WriteJSON(w, status, response)
}

func JSONErrorMessage(w http.ResponseWriter, status int, code string, message string) {
	response := ErrorMessageResponse{
		Success:   false,
		Code:      code,
		Details:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	WriteJSON(w, status, response)
}

func JSONErrorValidation(w http.ResponseWriter, err validator.ValidationErrorsTranslations) {
	resp := ErrorObjectResponse{
		Success:   false,
		Code:      "VALIDATION_FAILED",
		Message:   "Your submission contains some errors. Kindly review and fix the highlighted fields before submitting again.",
		Errors:    err,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	WriteJSON(w, http.StatusUnprocessableEntity, resp)
}

func Validate[T any](w http.ResponseWriter, form T) bool {
	// Initialize the English locale and translator
	english := en.New()
	uni := ut.New(english, english)
	translator, _ := uni.GetTranslator("en")

	// Register the default "en" translations to the validator
	validate := validator.New()

	// This changes the Field name in validation errors from "Full_Name" to "full_name"
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := entranslations.RegisterDefaultTranslations(validate, translator)
	if err != nil {
		JSONErrorMessage(w, http.StatusInternalServerError, "SERVER_ERROR", "Unable to register validation translator")
		return false
	}

	err = validate.Struct(form)
	if err != nil {
		var vErr validator.ValidationErrors

		// prepare a map to hold the validation errors
		cleanErrors := make(map[string]string)
		if errors.As(err, &vErr) {
			for _, f := range vErr {
				// f.Field() now returns the JSON name because of RegisterTagNameFunc
				cleanErrors[f.Field()] = f.Translate(translator)
			}
		}

		JSONErrorValidation(w, cleanErrors)
		return false
	}

	return true
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(true)

	if err := encoder.Encode(payload); err != nil {
		// last-resort fallback (never panic in HTTP)
		log.Printf("json encode error: %v", err)
		http.Error(w, `{"success":false, "message":"Whoops something went wrong"}`, http.StatusInternalServerError)
	}
}

func WriteEmpty(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
