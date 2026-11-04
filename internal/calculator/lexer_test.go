package calculator

import (
	"testing"
)

func TestLexerBasicTokens(t *testing.T) {
	input := "12 + 34.5 - 6 * 7 / 8 ^ 2"
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("unexpected error tokenizing: %v", err)
	}

	expectedTypes := []TokenType{
		TokenNumber, TokenPlus, TokenNumber, TokenMinus,
		TokenNumber, TokenMultiply, TokenNumber, TokenDivide,
		TokenNumber, TokenPower, TokenNumber,
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
		expectedTypes []TokenType
	}{
		{
			name:          "Leading unary minus",
			input:         "-5 + 3",
			expectedTypes: []TokenType{TokenUnaryMinus, TokenNumber, TokenPlus, TokenNumber},
		},
		{
			name:          "Unary minus after operator",
			input:         "4 * -2",
			expectedTypes: []TokenType{TokenNumber, TokenMultiply, TokenUnaryMinus, TokenNumber},
		},
		{
			name:          "Unary minus inside parens",
			input:         "(-3 + 1)",
			expectedTypes: []TokenType{TokenLParen, TokenUnaryMinus, TokenNumber, TokenPlus, TokenNumber, TokenRParen},
		},
		{
			name:          "Double minus (subtraction then negation)",
			input:         "5 - -2",
			expectedTypes: []TokenType{TokenNumber, TokenMinus, TokenUnaryMinus, TokenNumber},
		},
		{
			name:          "Unary plus",
			input:         "+7",
			expectedTypes: []TokenType{TokenUnaryPlus, TokenNumber},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lexer := NewLexer(tc.input)
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
		expectedTypes []TokenType
	}{
		{
			input:         "2(3)",
			expectedTypes: []TokenType{TokenNumber, TokenMultiply, TokenLParen, TokenNumber, TokenRParen},
		},
		{
			input:         "(2)(3)",
			expectedTypes: []TokenType{TokenLParen, TokenNumber, TokenRParen, TokenMultiply, TokenLParen, TokenNumber, TokenRParen},
		},
		{
			input:         "5sqrt(4)",
			expectedTypes: []TokenType{TokenNumber, TokenMultiply, TokenSqrt, TokenLParen, TokenNumber, TokenRParen},
		},
	}

	for _, tc := range tests {
		lexer := NewLexer(tc.input)
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
		{"Empty string", "", ErrCodeEmptyExpr},
		{"Only whitespace", "   \t\n  ", ErrCodeEmptyExpr},
		{"Multiple decimal points", "12.3.4", ErrCodeInvalidSyntax},
		{"Unexpected character", "5 + $", ErrCodeInvalidSyntax},
		{"Invalid percentage position", "% 5", ErrCodeInvalidSyntax},
		{"Unknown identifier", "foo(5)", ErrCodeInvalidSyntax},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lexer := NewLexer(tc.input)
			_, err := lexer.Tokenize()
			if err == nil {
				t.Fatalf("expected error for input %q, got nil", tc.input)
			}
			calcErr, ok := err.(*CalculationError)
			if !ok {
				t.Fatalf("expected *CalculationError, got %T", err)
			}
			if calcErr.Code != tc.errorCode {
				t.Errorf("expected error code %s, got %s", tc.errorCode, calcErr.Code)
			}
		})
	}
}
