package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/skb1129/calculator-app/internal/api"
	"github.com/skb1129/calculator-app/internal/api/handlers"
	"github.com/skb1129/calculator-app/internal/calculator"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestRouter() *gin.Engine {
	return api.SetupRouter("")
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

	// 10 + 5*8 - 8 + 20% = 10 + 40 - 8 + 20% = 42 + (42 * 0.20) = 42 + 8.4 = 50.4
	if res.Result != 50.4 {
		t.Errorf("expected result 50.4, got %v", res.Result)
	}
	if res.Formatted != "50.4" {
		t.Errorf("expected formatted '50.4', got %s", res.Formatted)
	}
}

func TestCalculateHandlerMaxLength(t *testing.T) {
	router := setupTestRouter()

	// Create an expression that exceeds MaxExpressionLength (1000 characters)
	longExpr := strings.Repeat("1+", 500) + "1" // 1001 characters
	if len(longExpr) <= handlers.MaxExpressionLength {
		t.Fatalf("longExpr length %d must exceed %d", len(longExpr), handlers.MaxExpressionLength)
	}

	payload := map[string]string{
		"expression": longExpr,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for too long expression, got %d: %s", w.Code, w.Body.String())
	}

	var errRes handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errRes); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if errRes.Code != handlers.ErrCodeExpressionTooLong {
		t.Errorf("expected code %s, got %s", handlers.ErrCodeExpressionTooLong, errRes.Code)
	}
	if !strings.Contains(errRes.Error, "exceeds maximum allowed length") {
		t.Errorf("expected error message to mention maximum length, got %q", errRes.Error)
	}
}

func TestCalculateHandlerRegexValidation(t *testing.T) {
	router := setupTestRouter()

	invalidExpressions := []struct {
		name string
		expr string
	}{
		{"Disallowed symbol dollar", "10 + $5"},
		{"Disallowed keyword alert", "alert(1)"},
		{"SQL injection attempt", "select * from users"},
		{"Unknown identifier foo", "foo(5)"},
		{"Disallowed variable x", "x + y"},
		{"OS command injection", "system('ls')"},
		{"Backticks or quotes", "5 + `3`"},
		{"HTML tag", "<script>1+1</script>"},
		{"Semicolon instruction", "5 + 3; drop table"},
	}

	for _, tc := range invalidExpressions {
		t.Run(tc.name, func(t *testing.T) {
			payload := map[string]string{
				"expression": tc.expr,
			}
			body, _ := json.Marshal(payload)

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400 for %q, got %d: %s", tc.expr, w.Code, w.Body.String())
			}

			var errRes handlers.ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &errRes); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if errRes.Code != calculator.ErrCodeInvalidSyntax {
				t.Errorf("expected code %s, got %s", calculator.ErrCodeInvalidSyntax, errRes.Code)
			}
			if !strings.Contains(errRes.Error, "invalid characters or unsupported keywords") {
				t.Errorf("expected invalid characters error, got %q", errRes.Error)
			}
		})
	}
}

func TestCalculateHandlerValidFormats(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name     string
		expr     string
		expected float64
	}{
		{"Unicode operators", "6 × 7 ÷ 2 + √(49)", 28},
		{"Scientific notation", "1e3 + 50", 1050},
		{"Nested functions and percentages", "sqrt(100) + 50%", 15},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payload := map[string]string{
				"expression": tc.expr,
			}
			body, _ := json.Marshal(payload)

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("expected status 200 for %q, got %d: %s", tc.expr, w.Code, w.Body.String())
			}

			var res handlers.CalculateResponse
			if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if res.Result != tc.expected {
				t.Errorf("expected result %v, got %v", tc.expected, res.Result)
			}
		})
	}
}

func TestCalculateHandlerEmptyAndWhitespace(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name string
		expr string
	}{
		{"Whitespace only", "    \t   "},
		{"Single space", " "},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payload := map[string]string{
				"expression": tc.expr,
			}
			body, _ := json.Marshal(payload)

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d", w.Code)
			}

			var errRes handlers.ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &errRes); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if errRes.Code != calculator.ErrCodeEmptyExpr {
				t.Errorf("expected code %s, got %s", calculator.ErrCodeEmptyExpr, errRes.Code)
			}
		})
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
