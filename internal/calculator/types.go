package calculator

// TokenType defines the kind of mathematical token.
type TokenType string

const (
	TokenNumber     TokenType = "NUMBER"
	TokenPlus       TokenType = "PLUS"
	TokenMinus      TokenType = "MINUS"
	TokenMultiply   TokenType = "MULTIPLY"
	TokenDivide     TokenType = "DIVIDE"
	TokenPower      TokenType = "POWER"
	TokenPercentage TokenType = "PERCENTAGE"
	TokenSqrt       TokenType = "SQRT"
	TokenLParen     TokenType = "LPAREN"
	TokenRParen     TokenType = "RPAREN"
	TokenUnaryMinus TokenType = "UNARY_MINUS"
	TokenUnaryPlus  TokenType = "UNARY_PLUS"
)

// Token represents a lexical token in a mathematical expression.
type Token struct {
	Type     TokenType
	Value    string
	Number   float64
	Position int
}

// CalculationError represents a domain or syntax error with contextual details.
type CalculationError struct {
	Message  string `json:"error"`
	Code     string `json:"code"`
	Position int    `json:"position,omitempty"`
}

func (e *CalculationError) Error() string {
	return e.Message
}

const (
	ErrCodeInvalidSyntax  = "INVALID_SYNTAX"
	ErrCodeEmptyExpr      = "EMPTY_EXPRESSION"
	ErrCodeUnmatchedParen = "UNMATCHED_PARENTHESIS"
	ErrCodeDivisionByZero = "DIVISION_BY_ZERO"
	ErrCodeMathDomain     = "MATH_DOMAIN_ERROR"
	ErrCodeMathOverflow   = "MATH_OVERFLOW"
)
