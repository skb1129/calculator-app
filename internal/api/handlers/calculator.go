package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skb1129/calculator-app/internal/calculator"
)

// MaxExpressionLength defines the maximum allowed character length for an expression string.
const MaxExpressionLength = 1000

// ErrCodeExpressionTooLong is returned when an input expression exceeds MaxExpressionLength.
const ErrCodeExpressionTooLong = "EXPRESSION_TOO_LONG"

// ValidExpressionRegex enforces that expressions contain only supported mathematical tokens, numbers, and functions.
var ValidExpressionRegex = regexp.MustCompile(`(?i)^([0-9\s+\-*/^%().×÷√]|sqrt|[0-9]+(\.[0-9]+)?e[+-]?[0-9]+)+$`)

// CalculateRequest models incoming expression payload.
type CalculateRequest struct {
	Expression string `json:"expression" binding:"required"`
}

// CalculateResponse represents successful calculation response.
type CalculateResponse struct {
	Success         bool     `json:"success"`
	Expression      string   `json:"expression"`
	Result          float64  `json:"result"`
	Formatted       string   `json:"formatted"`
	Tokens          []string `json:"tokens,omitempty"`
	ExecutionTimeMs float64  `json:"execution_time_ms"`
}

// ErrorResponse represents structured error details.
type ErrorResponse struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	Code     string `json:"code"`
	Position int    `json:"position,omitempty"`
}

// Calculate handles POST /api/v1/calculate.
func Calculate(c *gin.Context) {
	var req CalculateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Expression field is required and must be a valid JSON string",
			Code:    calculator.ErrCodeInvalidSyntax,
		})
		return
	}

	trimmed := strings.TrimSpace(req.Expression)
	if len(trimmed) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Expression cannot be empty",
			Code:    calculator.ErrCodeEmptyExpr,
		})
		return
	}

	// 1. Max length validation
	if len(req.Expression) > MaxExpressionLength {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   fmt.Sprintf("Expression exceeds maximum allowed length of %d characters", MaxExpressionLength),
			Code:    ErrCodeExpressionTooLong,
		})
		return
	}

	// 2. Regex-based pattern validation
	if !ValidExpressionRegex.MatchString(trimmed) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Expression contains invalid characters or unsupported keywords",
			Code:    calculator.ErrCodeInvalidSyntax,
		})
		return
	}

	start := time.Now()
	res, err := calculator.Calculate(req.Expression)
	durationMs := float64(time.Since(start).Nanoseconds()) / 1e6

	if err != nil {
		if calcErr, ok := errors.AsType[*calculator.CalculationError](err); ok {
			statusCode := http.StatusBadRequest
			if calcErr.Code == calculator.ErrCodeDivisionByZero ||
				calcErr.Code == calculator.ErrCodeMathDomain ||
				calcErr.Code == calculator.ErrCodeMathOverflow {
				statusCode = http.StatusUnprocessableEntity
			}

			c.JSON(statusCode, ErrorResponse{
				Success:  false,
				Error:    calcErr.Message,
				Code:     calcErr.Code,
				Position: calcErr.Position,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Internal error evaluating expression",
			Code:    "INTERNAL_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, CalculateResponse{
		Success:         true,
		Expression:      res.Expression,
		Result:          res.Result,
		Formatted:       res.Formatted,
		Tokens:          res.Tokens,
		ExecutionTimeMs: durationMs,
	})
}
