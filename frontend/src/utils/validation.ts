export interface ValidationResult {
  isValid: boolean;
  error?: string;
  position?: number;
}

export function validateExpression(expr: string): ValidationResult {
  const trimmed = expr.trim();
  if (!trimmed) {
    return {
      isValid: false,
      error: 'Please enter an expression to calculate',
    };
  }

  // Bracket balance check
  let balance = 0;
  for (let i = 0; i < trimmed.length; i++) {
    const char = trimmed[i];
    if (char === '(') balance++;
    if (char === ')') {
      balance--;
      if (balance < 0) {
        return {
          isValid: false,
          error: `Unmatched closing parenthesis ')' at character ${i + 1}`,
          position: i,
        };
      }
    }
  }

  if (balance > 0) {
    return {
      isValid: false,
      error: `Missing ${balance} closing parenthesis ')'`,
    };
  }

  // Check for consecutive binary operators (e.g., "+ *", "* /", "/ +")
  const invalidOperatorSeq = /[+\-*/^×÷]{2,}/;
  const match = invalidOperatorSeq.exec(trimmed.replace(/\s+/g, ''));
  if (match) {
    // Note: unary minus after operator like "*-" or "+-" is valid math!
    const seq = match[0];
    if (!seq.endsWith('-') || seq.length > 2) {
      return {
        isValid: false,
        error: `Invalid consecutive operators: '${seq}'`,
      };
    }
  }

  // Check for multiple decimal points in a single number token
  const tokens = trimmed.split(/[^0-9.]+/);
  for (const token of tokens) {
    if (token.split('.').length > 2) {
      return {
        isValid: false,
        error: `Invalid number format with multiple decimal points: '${token}'`,
      };
    }
  }

  return { isValid: true };
}
