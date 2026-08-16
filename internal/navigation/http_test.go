package navigation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResourceImportEndpointReturnsImportedCount(t *testing.T) {
	service := NewServiceWithFixtures()
	handler := NewHandler(service)
	payload := []byte(`{"links":[{"id":"library","group":"教师工具","title":"数字图书馆","url":"https://resources.example.test/library","sortOrder":15}]}`)
	request := httptest.NewRequest(http.MethodPost, "/api/resources/import", bytes.NewReader(payload))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	var result ImportResult
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result.Imported != 1 || result.Rejected != 0 {
		t.Fatalf("unexpected import response: %+v", result)
	}
	linksRequest := httptest.NewRequest(http.MethodGet, "/api/resources?group=%E6%95%99%E5%B8%88%E5%B7%A5%E5%85%B7", nil)
	linksResponse := httptest.NewRecorder()
	handler.ServeHTTP(linksResponse, linksRequest)
	if linksResponse.Code != http.StatusOK {
		t.Fatalf("expected query status 200, got %d", linksResponse.Code)
	}
}
