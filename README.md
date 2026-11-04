# OmniCalc: Full-Stack Multi-Operation Calculator

OmniCalc is a production-grade full-stack calculator application capable of parsing, validating, and evaluating complex, multi-operand mathematical expressions server-side with high numerical precision.

- **Frontend**: Responsive Single Page Application (SPA) built with **React 19**, **TypeScript**, and **Tailwind CSS**.
- **Backend**: High-performance REST service built in **Go** using the **Gin** web framework.
- **Parsing Engine**: Custom mathematical parser implementing Dijkstra's **Shunting-Yard algorithm** and Reverse Polish Notation (RPN) evaluation.
- **Deployment**: Multi-stage **Dockerfile** and **Docker Compose** packaging the frontend and backend into a single lightweight container (< 25MB).

---

## Table of Contents
- [Features](#features)
- [Architecture & Design Rationale](#architecture--design-rationale)
- [Supported Mathematical Operations](#supported-mathematical-operations)
- [API Reference & Usage](#api-reference--usage)
- [Project Structure](#project-structure)
- [Local Development & Setup](#local-development--setup)
- [Testing](#testing)
- [Docker Deployment](#docker-deployment)
- [License](#license)

---

## Features

- **Complex Expression Evaluation**: Handles multi-operand, multi-operator strings (e.g., `(15 + 5 * 2^3) / sqrt(16) - 10%`).
- **Mathematical Integrity**: Respects strict operator precedence, right-associative exponentiation ($2^{3^2} = 512$), prefix unary minus (`-5 + 3`, `3 * -2`), and postfix percentages.
- **In-Memory History Log**: Tracks calculation history with timestamps and allows one-click loading of past expressions or inserting results into active calculations.
- **Responsive & Accessible Design**: Intuitive keypad, direct text input, physical keyboard bindings, tactile feedback, and seamless layout adaptation across mobile, tablet, and desktop viewports.
- **Comprehensive Error Handling**: Validates expressions client-side before dispatch and surfaces structured server-side errors (division by zero, negative square root, syntax violations).
- **Floating-Point Precision Formatting**: Eliminates IEEE 754 precision artifacts (e.g., formats `0.1 + 0.2` cleanly as `0.3`).

---

## Architecture & Design Rationale

### 1. Server-Side Expression Parsing
Rather than evaluating calculations client-side with insecure JavaScript `eval()` or rudimentary token splitting, all calculations are submitted to a dedicated Go service. Go's compiled performance, minimal memory footprint, and strong type safety make it ideal for mathematical computation.

### 2. Lexer & Shunting-Yard Parser (RPN)
The backend evaluates expressions using a two-phase architecture:
1. **Lexical Analysis (`lexer.go`)**: Scans character streams into structured tokens, automatically distinguishing binary subtraction from unary negation based on token context. Supports implicit multiplication (e.g., `2(3)` $\to$ `2 * 3`, `5sqrt(4)` $\to$ `5 * sqrt(4)`).
2. **Shunting-Yard Algorithm (`parser.go`)**: Converts infix token streams into Reverse Polish Notation (RPN), resolving operator precedence and associativity without recursive AST overhead.
3. **RPN Evaluator (`evaluator.go`)**: Employs an operand stack to compute results in $O(n)$ time while checking mathematical invariants (e.g., zero divisors, negative square root domains).

### 3. Operator Semantics
- **Percentage (`%`)**: Modeled as a postfix unary operator ($x\% = x / 100$). This ensures mathematical consistency across composite equations: `100 * 15%` evaluates to `15`, and `50 + 10%` evaluates to `50.1`.
- **Exponentiation (`^`)**: Right-associative per standard mathematical conventions ($2^{3^2} = 2^{(3^2)} = 512$).
- **Square Root (`sqrt` / `√`)**: Supported both as functional calls (`sqrt(16)`) and with Unicode symbols (`√(16)` or `√16`).

### 4. Unified Full-Stack Container
In production, a multi-stage `Dockerfile` compiles the React application into static assets and builds the Go binary. The Go Gin service statically serves the compiled frontend assets while exposing the `/api/v1` routes, eliminating CORS overhead in production and removing the need for a separate Nginx container.

---

## Supported Mathematical Operations

| Operation | Operator Syntax | Precedence | Associativity | Example | Notes |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Parentheses** | `(`, `)` | 0 (Group) | N/A | `(2 + 3) * 4` | Sub-expression isolation |
| **Percentage** | `%` | 5 (Postfix) | Left | `50 * 20%` | Evaluates to $x / 100$ |
| **Square Root** | `sqrt(x)` or `√(x)` | 5 (Prefix) | Right | `sqrt(144) + 1` | Requires $x \ge 0$ |
| **Exponentiation** | `^` | 4 (Binary) | Right | `2 ^ 3 ^ 2` | Evaluates right-to-left |
| **Unary Negation** | `-` | 3 (Prefix) | Right | `-5 + 3`, `4 * -2` | Precedence above binary operators |
| **Unary Plus** | `+` | 3 (Prefix) | Right | `+7` | Identity operation |
| **Multiplication** | `*` or `×` | 2 (Binary) | Left | `6 * 7`, `6 × 7` | Standard binary multiplication |
| **Division** | `/` or `÷` | 2 (Binary) | Left | `20 / 4`, `20 ÷ 4` | Guarded against division by zero |
| **Addition** | `+` | 1 (Binary) | Left | `12 + 34` | Standard binary addition |
| **Subtraction** | `-` | 1 (Binary) | Left | `45 - 12` | Standard binary subtraction |

---

## API Reference & Usage

### 1. Calculate Expression
**Endpoint:** `POST /api/v1/calculate`  
**Content-Type:** `application/json`

#### Request Body
```json
{
  "expression": "(15 + 5 * 2^3) / sqrt(16) - 10%"
}
```

#### Success Response (`200 OK`)
```json
{
  "success": true,
  "expression": "(15 + 5 * 2^3) / sqrt(16) - 10%",
  "result": 13.65,
  "formatted": "13.65",
  "tokens": ["(", "15", "+", "5", "*", "2", "^", "3", ")", "/", "sqrt", "(", "16", ")", "-", "10", "%"],
  "execution_time_ms": 0.08
}
```

#### Error Response (`400 Bad Request` or `422 Unprocessable Entity`)
```json
{
  "success": false,
  "error": "Division by zero at position 12",
  "code": "DIVISION_BY_ZERO",
  "position": 12
}
```

#### Error Codes:
- `EMPTY_EXPRESSION`: Expression string was empty or whitespace.
- `INVALID_SYNTAX`: Malformed numbers, invalid tokens, or consecutive invalid operators.
- `UNMATCHED_PARENTHESIS`: Mismatched opening `(` or closing `)` bracket.
- `DIVISION_BY_ZERO`: Division by zero ($x / 0$).
- `MATH_DOMAIN_ERROR`: Real-number domain violations (e.g., $\sqrt{-9}$ or negative numbers with fractional powers).
- `MATH_OVERFLOW`: Numeric overflow outside 64-bit float range.

#### Example `curl`:
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression": "10 + 5 * 2^3 - sqrt(64) + 20%"}'
```

---

### 2. Health Check
**Endpoint:** `GET /api/v1/health`

#### Response (`200 OK`)
```json
{
  "status": "healthy",
  "timestamp": "2026-11-04T08:30:00Z",
  "service": "calculator-service"
}
```

#### Example `curl`:
```bash
curl http://localhost:8080/api/v1/health
```

---

## Project Structure

```
calculator-app/
├── cmd/
│   └── server/
│       └── main.go                 # Gin server entrypoint & router initialization
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── calculator.go       # POST /api/v1/calculate handler
│   │   │   ├── calculator_test.go  # Handlers integration test suite
│   │   │   └── health.go           # GET /api/v1/health handler
│   │   ├── middleware/
│   │   │   └── cors.go             # CORS configuration middleware
│   │   └── routes.go               # Router registry & static SPA routing
│   └── calculator/
│       ├── evaluator.go            # RPN stack evaluator & precision formatter
│       ├── evaluator_test.go       # Mathematical engine unit test suite
│       ├── lexer.go                # Tokenizer & unary context resolver
│       ├── lexer_test.go           # Lexical analysis unit test suite
│       ├── parser.go               # Shunting-yard infix-to-RPN parser
│       └── types.go                # Token definitions & error types
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Calculator.tsx      # Main calculator layout
│   │   │   ├── Display.tsx         # Expression input & result display
│   │   │   ├── ErrorAlert.tsx      # Structured error notification banner
│   │   │   ├── HistoryLog.tsx      # In-memory history list with recall
│   │   │   └── Keypad.tsx          # Responsive tactile button grid
│   │   ├── hooks/
│   │   │   └── useCalculator.ts    # State management & keyboard listener
│   │   ├── services/
│   │   │   └── api.ts              # Backend API client & health check
│   │   ├── types/
│   │   │   └── calculator.ts       # Shared TypeScript interfaces
│   │   ├── utils/
│   │   │   ├── formatting.ts       # Number & timestamp formatting
│   │   │   └── validation.ts       # Pre-dispatch syntax validation
│   │   ├── App.tsx                 # Root application wrapper
│   │   ├── main.tsx                # React DOM entry point
│   │   └── index.css               # Tailwind CSS declarations
│   ├── tests/
│   │   ├── Calculator.test.tsx     # React component integration tests
│   │   ├── setup.ts                # Test environment setup
│   │   └── validation.test.ts      # Client validation unit tests
│   ├── package.json
│   ├── tsconfig.json
│   ├── tailwind.config.js
│   └── vite.config.ts
├── Dockerfile                      # Production multi-stage Dockerfile
├── docker-compose.yml              # Local container orchestration
├── go.mod
├── go.sum
└── README.md
```

---

## Local Development & Setup

### Prerequisites
- **Go**: 1.22+ (tested on Go 1.24/1.27)
- **Node.js**: v20+ or v22+
- **npm**: v10+

### 1. Run Backend Service
```bash
# From repository root
go run cmd/server/main.go
```
The server will start on port `8080`.

### 2. Run Frontend Development Server
In a separate terminal:
```bash
cd frontend
npm install
npm run dev
```
Open [http://localhost:3000](http://localhost:3000) in your browser. The Vite development server automatically proxies API requests to `http://localhost:8080`.

---

## Testing

### Run Backend Unit Tests (Go)
```bash
go test ./... -v -cover
```
Runs the full suite of tokenizer, parser, evaluator, and HTTP integration tests with statement coverage metrics.

### Run Frontend Unit Tests (Vitest + React Testing Library)
```bash
cd frontend
npm run test
```
Executes component integration tests, API mock tests, keyboard handler checks, and syntax validation tests.

---

## Docker Deployment

### 1. Build and Run with Docker Compose
The simplest way to build and run the full-stack application is with Docker Compose:

```bash
docker compose up --build
```
The application will be accessible at [http://localhost:8080](http://localhost:8080).

### 2. Build and Run Standalone Docker Image
```bash
# Build multi-stage image
docker build -t calculator-app:latest .

# Run container
docker run -p 8080:8080 --name calculator calculator-app:latest
```

The multi-stage build automatically:
1. Builds the React SPA (`npm run build`).
2. Statically compiles the Go server binary without external C dependencies (`CGO_ENABLED=0`).
3. Copies both into a minimal Alpine Linux container running under an unprivileged user (`appuser`).
4. Serves both the REST API and the frontend assets from a single container on port `8080`.

---

## Keyboard Shortcuts

| Key | Action |
| :--- | :--- |
| `0` - `9`, `.` | Enter numbers and decimals |
| `+`, `-`, `*`, `/` | Enter binary operations |
| `^` | Enter exponentiation |
| `%` | Enter percentage |
| `(`, `)` | Enter parentheses |
| `Enter` or `=` | Submit calculation to backend |
| `Backspace` | Delete previous character/token |
| `Escape` | Clear current expression (`AC`) |
