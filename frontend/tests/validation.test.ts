import { describe, it, expect } from 'vitest';
import { validateExpression } from '../src/utils/validation';

describe('validateExpression', () => {
  it('should accept valid multi-operation expressions', () => {
    expect(validateExpression('2 + 3 * 4').isValid).toBe(true);
    expect(validateExpression('(10 + 5) / 2').isValid).toBe(true);
    expect(validateExpression('sqrt(16) + 2^3 - 50%').isValid).toBe(true);
    expect(validateExpression('4 * -2').isValid).toBe(true);
  });

  it('should reject empty or whitespace-only expressions', () => {
    expect(validateExpression('').isValid).toBe(false);
    expect(validateExpression('   ').isValid).toBe(false);
  });

  it('should reject unmatched parentheses', () => {
    const unclosed = validateExpression('(2 + 3 * 4');
    expect(unclosed.isValid).toBe(false);
    expect(unclosed.error).toContain('Missing 1 closing parenthesis');

    const extraClosing = validateExpression('2 + 3) * 4');
    expect(extraClosing.isValid).toBe(false);
    expect(extraClosing.error).toContain("Unmatched closing parenthesis ')'");
  });

  it('should reject multiple decimal points in a single number', () => {
    const result = validateExpression('12.3.4 + 5');
    expect(result.isValid).toBe(false);
    expect(result.error).toContain('multiple decimal points');
  });

  it('should reject invalid consecutive operators', () => {
    const result = validateExpression('5 ++ 2');
    expect(result.isValid).toBe(false);
    expect(result.error).toContain('Invalid consecutive operators');
  });
});
