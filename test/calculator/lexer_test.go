package calculator_test

import (
	"testing"

	"github.com/skb1129/calculator-app/internal/calculator"
)

func TestLexerBasicTokens(t *testing.T) {
	input := "12 + 34.5 - 6 * 7 / 8 ^ 2"
	lexer := calculator.NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("unexpected error tokenizing: %v", err)
	}

	expectedTypes := []calculator.TokenType{
		calculator.TokenNumber, calculator.TokenPlus, calculator.TokenNumber, calculator.TokenMinus,
		calculator.TokenNumber, calculator.TokenMultiply, calculator.TokenNumber, calculator.TokenDivide,
		calculator.TokenNumber, calculator.TokenPower, calculator.TokenNumber,
	}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, expected := range expectedTypes {
		if tokens[i].Type != expected {
			t.Errorf("token %d: expected type %s, got %s", i, expected, tokens[i].Type)
		}
	}
}

func TestLexerUnaryMinusAndPlus(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedTypes []calculator.TokenType
	}{
		{
			name:          "Leading unary minus",
			input:         "-5 + 3",
			expectedTypes: []calculator.TokenType{calculator.TokenUnaryMinus, calculator.TokenNumber, calculator.TokenPlus, calculator.TokenNumber},
		},
		{
			name:          "Unary minus after operator",
			input:         "4 * -2",
			expectedTypes: []calculator.TokenType{calculator.TokenNumber, calculator.TokenMultiply, calculator.TokenUnaryMinus, calculator.TokenNumber},
		},
		{
			name:          "Unary minus inside parens",
			input:         "(-3 + 1)",
			expectedTypes: []calculator.TokenType{calculator.TokenLParen, calculator.TokenUnaryMinus, calculator.TokenNumber, calculator.TokenPlus, calculator.TokenNumber, calculator.TokenRParen},
		},
		{
			name:          "Double minus (subtraction then negation)",
			input:         "5 - -2",
			expectedTypes: []calculator.TokenType{calculator.TokenNumber, calculator.TokenMinus, calculator.TokenUnaryMinus, calculator.TokenNumber},
		},
		{
			name:          "Unary plus",
			input:         "+7",
			expectedTypes: []calculator.TokenType{calculator.TokenUnaryPlus, calculator.TokenNumber},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lexer := calculator.NewLexer(tc.input)
			tokens, err := lexer.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tokens) != len(tc.expectedTypes) {
				t.Fatalf("expected %d tokens, got %d", len(tc.expectedTypes), len(tokens))
			}
			for i, exp := range tc.expectedTypes {
				if tokens[i].Type != exp {
					t.Errorf("token %d: expected %s, got %s", i, exp, tokens[i].Type)
				}
			}
		})
	}
}

func TestLexerImplicitMultiplication(t *testing.T) {
	tests := []struct {
		input         string
		expectedTypes []calculator.TokenType
	}{
		{
			input:         "2(3)",
			expectedTypes: []calculator.TokenType{calculator.TokenNumber, calculator.TokenMultiply, calculator.TokenLParen, calculator.TokenNumber, calculator.TokenRParen},
		},
		{
			input:         "(2)(3)",
			expectedTypes: []calculator.TokenType{calculator.TokenLParen, calculator.TokenNumber, calculator.TokenRParen, calculator.TokenMultiply, calculator.TokenLParen, calculator.TokenNumber, calculator.TokenRParen},
		},
		{
			input:         "5sqrt(4)",
			expectedTypes: []calculator.TokenType{calculator.TokenNumber, calculator.TokenMultiply, calculator.TokenSqrt, calculator.TokenLParen, calculator.TokenNumber, calculator.TokenRParen},
		},
	}

	for _, tc := range tests {
		lexer := calculator.NewLexer(tc.input)
		tokens, err := lexer.Tokenize()
		if err != nil {
			t.Fatalf("failed for input %s: %v", tc.input, err)
		}
		if len(tokens) != len(tc.expectedTypes) {
			t.Fatalf("input %s: expected %d tokens, got %d", tc.input, len(tc.expectedTypes), len(tokens))
		}
		for i, exp := range tc.expectedTypes {
			if tokens[i].Type != exp {
				t.Errorf("input %s token %d: expected %s, got %s", tc.input, i, exp, tokens[i].Type)
			}
		}
	}
}

func TestLexerErrors(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		errorCode string
	}{
		{"Empty string", "", calculator.ErrCodeEmptyExpr},
		{"Only whitespace", "   \t\n  ", calculator.ErrCodeEmptyExpr},
		{"Multiple decimal points", "12.3.4", calculator.ErrCodeInvalidSyntax},
		{"Unexpected character", "5 + $", calculator.ErrCodeInvalidSyntax},
		{"Invalid percentage position", "% 5", calculator.ErrCodeInvalidSyntax},
		{"Unknown identifier", "foo(5)", calculator.ErrCodeInvalidSyntax},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lexer := calculator.NewLexer(tc.input)
			_, err := lexer.Tokenize()
			if err == nil {
				t.Fatalf("expected error for input %q, got nil", tc.input)
			}
			calcErr, ok := err.(*calculator.CalculationError)
			if !ok {
				t.Fatalf("expected *calculator.CalculationError, got %T", err)
			}
			if calcErr.Code != tc.errorCode {
				t.Errorf("expected error code %s, got %s", tc.errorCode, calcErr.Code)
			}
		})
	}
}
