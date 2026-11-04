package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/skb1129/calculator-app/internal/api"
	"github.com/skb1129/calculator-app/internal/api/handlers"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestRouter() *gin.Engine {
	return api.SetupRouter("")
}

func TestHealthCheckHandler(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var res handlers.HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Status != "healthy" {
		t.Errorf("expected status healthy, got %s", res.Status)
	}
}

func TestCalculateHandlerSuccess(t *testing.T) {
	router := setupTestRouter()

	payload := map[string]string{
		"expression": "10 + 5 * 2^3 - sqrt(64) + 20%",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var res handlers.CalculateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success true, got false")
	}

	// 10 + 5*8 - 8 + 0.2 = 10 + 40 - 8 + 0.2 = 42.2
	if res.Result != 42.2 {
		t.Errorf("expected result 42.2, got %v", res.Result)
	}
	if res.Formatted != "42.2" {
		t.Errorf("expected formatted '42.2', got %s", res.Formatted)
	}
}

func TestCalculateHandlerSyntaxError(t *testing.T) {
	router := setupTestRouter()

	payload := map[string]string{
		"expression": "(2 + 3 * 4",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}

	var errRes handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errRes); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if errRes.Success {
		t.Errorf("expected success false, got true")
	}
	if errRes.Code != "UNMATCHED_PARENTHESIS" {
		t.Errorf("expected code UNMATCHED_PARENTHESIS, got %s", errRes.Code)
	}
}

func TestCalculateHandlerDivisionByZero(t *testing.T) {
	router := setupTestRouter()

	payload := map[string]string{
		"expression": "42 / (5 - 5)",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d: %s", w.Code, w.Body.String())
	}

	var errRes handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errRes); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if errRes.Code != "DIVISION_BY_ZERO" {
		t.Errorf("expected code DIVISION_BY_ZERO, got %s", errRes.Code)
	}
}

func TestCalculateHandlerNegativeSqrt(t *testing.T) {
	router := setupTestRouter()

	payload := map[string]string{
		"expression": "sqrt(-25)",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d: %s", w.Code, w.Body.String())
	}

	var errRes handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errRes); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if errRes.Code != "MATH_DOMAIN_ERROR" {
		t.Errorf("expected code MATH_DOMAIN_ERROR, got %s", errRes.Code)
	}
}

func TestCalculateHandlerCORS(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest(http.MethodOptions, "/api/v1/calculate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 No Content for OPTIONS, got %d", w.Code)
	}

	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "*" {
		t.Errorf("expected Access-Control-Allow-Origin *, got %s", allowOrigin)
	}
}
