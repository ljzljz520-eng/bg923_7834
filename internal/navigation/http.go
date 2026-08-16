package navigation

import (
	"encoding/json"
	"net/http"
)

import "github.com/go-chi/chi/v5"

type importRequest struct {
	Links []ImportLink `json:"links"`
}

func NewHandler(service *Service) http.Handler {
	router := chi.NewRouter()
	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, http.StatusOK, service.Home())
	})
	router.Get("/api/home", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, http.StatusOK, service.Home())
	})
	router.Get("/api/groups", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, http.StatusOK, map[string]any{"groups": service.Groups()})
	})
	router.Get("/api/resources", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, http.StatusOK, map[string]any{"links": service.List(request.URL.Query().Get("group"))})
	})
	router.Get("/api/resources/status", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, http.StatusOK, map[string]any{"resources": service.Validate(request.URL.Query().Get("group"))})
	})
	router.Post("/api/resources/import", func(writer http.ResponseWriter, request *http.Request) {
		var payload importRequest
		decoder := json.NewDecoder(request.Body)
		if err := decoder.Decode(&payload); err != nil {
			writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "request body must be valid JSON"})
			return
		}
		result := service.ImportBatch(payload.Links)
		status := http.StatusOK
		if result.Imported == 0 && result.Rejected > 0 {
			status = http.StatusUnprocessableEntity
		}
		writeJSON(writer, status, result)
	})
	return router
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}
