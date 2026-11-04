package calculator

import (
	"math"
	"testing"
)

func TestCalculateBasicOperations(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		formatted  string
	}{
		{"2 + 3", 5, "5"},
		{"10 - 4", 6, "6"},
		{"6 * 7", 42, "42"},
		{"6 × 7", 42, "42"},
		{"20 / 4", 5, "5"},
		{"20 ÷ 4", 5, "5"},
		{"2 ^ 3", 8, "8"},
		{"sqrt(64)", 8, "8"},
		{"√(49)", 7, "7"},
		{"50%", 0.5, "0.5"},
	}

	for _, tc := range tests {
		t.Run(tc.expression, func(t *testing.T) {
			res, err := Calculate(tc.expression)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.expression, err)
			}
			if math.Abs(res.Result-tc.expected) > 1e-9 {
				t.Errorf("result: expected %v, got %v", tc.expected, res.Result)
			}
			if res.Formatted != tc.formatted {
				t.Errorf("formatted: expected %q, got %q", tc.formatted, res.Formatted)
			}
		})
	}
}

func TestCalculatePrecedenceAndAssociativity(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   float64
	}{
		{"Multiplication before addition", "2 + 3 * 4", 14},
		{"Division before subtraction", "20 - 10 / 2", 15},
		{"Parentheses override precedence", "(2 + 3) * 4", 20},
		{"Nested parentheses", "((2 + 3) * (4 - 1)) / 5", 3},
		{"Right associativity for power", "2 ^ 3 ^ 2", 512},     // 2^(3^2) = 2^9 = 512
		{"Explicit left grouping for power", "(2 ^ 3) ^ 2", 64}, // (2^3)^2 = 8^2 = 64
		{"Power before multiplication", "2 * 3 ^ 2", 18},        // 2 * 9 = 18
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Calculate(tc.expression)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(res.Result-tc.expected) > 1e-9 {
				t.Errorf("expected %v, got %v", tc.expected, res.Result)
			}
		})
	}
}

func TestCalculateUnaryMinus(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
	}{
		{"-5 + 3", -2},
		{"3 * -2", -6},
		{"-3 ^ 2", -9},    // -(3^2) = -9
		{"(-3) ^ 2", 9},   // (-3)^2 = 9
		{"2 ^ -3", 0.125}, // 2^(-3) = 0.125
		{"5 - -2", 7},
		{"-(4 + 3)", -7},
		{"--5", 5},
	}

	for _, tc := range tests {
		t.Run(tc.expression, func(t *testing.T) {
			res, err := Calculate(tc.expression)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.expression, err)
			}
			if math.Abs(res.Result-tc.expected) > 1e-9 {
				t.Errorf("for %s: expected %v, got %v", tc.expression, tc.expected, res.Result)
			}
		})
	}
}

func TestCalculatePercentage(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
	}{
		{"100 * 15%", 15},
		{"50 + 10%", 50.1},
		{"200 - 50%", 199.5},
		{"10% + 20%", 0.3},
		{"(100 + 20)%", 1.2},
	}

	for _, tc := range tests {
		t.Run(tc.expression, func(t *testing.T) {
			res, err := Calculate(tc.expression)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.expression, err)
			}
			if math.Abs(res.Result-tc.expected) > 1e-9 {
				t.Errorf("for %s: expected %v, got %v", tc.expression, tc.expected, res.Result)
			}
		})
	}
}

func TestCalculateSquareRoot(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
	}{
		{"sqrt(144) + 1", 13},
		{"sqrt(9 + 16)", 5},
		{"2 * sqrt(25)", 10},
		{"sqrt(sqrt(16))", 2},
		{"√(81) / 3", 3},
	}

	for _, tc := range tests {
		t.Run(tc.expression, func(t *testing.T) {
			res, err := Calculate(tc.expression)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.expression, err)
			}
			if math.Abs(res.Result-tc.expected) > 1e-9 {
				t.Errorf("for %s: expected %v, got %v", tc.expression, tc.expected, res.Result)
			}
		})
	}
}

func TestCalculateFloatingPointPrecision(t *testing.T) {
	res, err := Calculate("0.1 + 0.2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Formatted != "0.3" {
		t.Errorf("expected clean formatted float '0.3', got %q", res.Formatted)
	}
}

func TestCalculateComplexExpression(t *testing.T) {
	// (15 + 5 * 2^3) / sqrt(16) - 10%
	// = (15 + 5 * 8) / 4 - 0.1
	// = (15 + 40) / 4 - 0.1
	// = 55 / 4 - 0.1
	// = 13.75 - 0.1 = 13.65
	expr := "(15 + 5 * 2^3) / sqrt(16) - 10%"
	res, err := Calculate(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := 13.65
	if math.Abs(res.Result-expected) > 1e-9 {
		t.Errorf("expected %v, got %v", expected, res.Result)
	}
	if res.Formatted != "13.65" {
		t.Errorf("expected formatted '13.65', got %q", res.Formatted)
	}
}

func TestCalculateErrors(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		errorCode  string
	}{
		{"Division by zero", "10 / 0", ErrCodeDivisionByZero},
		{"Division by zero in subexpression", "5 / (3 - 3)", ErrCodeDivisionByZero},
		{"Square root of negative number", "sqrt(-9)", ErrCodeMathDomain},
		{"Square root of negative in subexpression", "sqrt(2 - 10)", ErrCodeMathDomain},
		{"Unmatched opening parenthesis", "(2 + 3", ErrCodeUnmatchedParen},
		{"Unmatched closing parenthesis", "2 + 3)", ErrCodeUnmatchedParen},
		{"Empty expression", "", ErrCodeEmptyExpr},
		{"Negative base to fractional power", "(-4) ^ 0.5", ErrCodeMathDomain},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Calculate(tc.expression)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", tc.expression)
			}
			calcErr, ok := err.(*CalculationError)
			if !ok {
				t.Fatalf("expected *CalculationError, got %T: %v", err, err)
			}
			if calcErr.Code != tc.errorCode {
				t.Errorf("expected error code %s, got %s (msg: %s)", tc.errorCode, calcErr.Code, calcErr.Message)
			}
		})
	}
}
