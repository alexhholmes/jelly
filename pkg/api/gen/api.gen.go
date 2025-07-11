//go:build go1.22

// Package gen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package gen

import (
	"fmt"
	"net/http"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// BadRequest defines model for BadRequest.
type BadRequest struct {
	Message string `json:"message"`
}

// Forbidden defines model for Forbidden.
type Forbidden struct {
	Message string `json:"message"`
}

// HealthCheck defines model for HealthCheck.
type HealthCheck struct {
	Status string `json:"status"`
}

// InternalServerError defines model for InternalServerError.
type InternalServerError struct {
	Message string `json:"message"`
}

// Photo defines model for Photo.
type Photo struct {
	// Caption Photo caption
	Caption *string `json:"caption,omitempty"`

	// Id Unique identifier for the photo
	Id string `json:"id"`

	// Tags Photo tags
	Tags *[]string `json:"tags,omitempty"`

	// UploadedAt Timestamp when photo was uploaded
	UploadedAt time.Time `json:"uploadedAt"`

	// Url URL to access the uploaded photo
	Url string `json:"url"`
}

// PhotoUploadResponse defines model for PhotoUploadResponse.
type PhotoUploadResponse struct {
	Message *string `json:"message,omitempty"`
	Photo   Photo   `json:"photo"`
}

// InternalError defines model for internal-error.
type InternalError = InternalServerError

// UploadPhotoMultipartBody defines parameters for UploadPhoto.
type UploadPhotoMultipartBody struct {
	// Caption Optional caption for the photo
	Caption *string `json:"caption,omitempty"`

	// File The photo file to upload
	File openapi_types.File `json:"file"`

	// Tags Optional tags for the photo
	Tags *[]string `json:"tags,omitempty"`
}

// UploadPhotoMultipartRequestBody defines body for UploadPhoto for multipart/form-data ContentType.
type UploadPhotoMultipartRequestBody UploadPhotoMultipartBody

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /health)
	HealthCheck(w http.ResponseWriter, r *http.Request)

	// (POST /photo)
	UploadPhoto(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// HealthCheck operation middleware
func (siw *ServerInterfaceWrapper) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.HealthCheck(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// UploadPhoto operation middleware
func (siw *ServerInterfaceWrapper) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UploadPhoto(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       *http.ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m *http.ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m *http.ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("GET "+options.BaseURL+"/health", wrapper.HealthCheck)
	m.HandleFunc("POST "+options.BaseURL+"/photo", wrapper.UploadPhoto)

	return m
}
