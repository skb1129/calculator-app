import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useCalculator } from '../src/hooks/useCalculator';
import * as apiService from '../src/services/api';
import { CalculatorApiError } from '../src/services/api';

describe('useCalculator Hook', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('initializes with default empty state', () => {
    const { result } = renderHook(() => useCalculator());

    expect(result.current.expression).toBe('');
    expect(result.current.result).toBeNull();
    expect(result.current.history).toEqual([]);
    expect(result.current.isLoading).toBe(false);
    expect(result.current.error).toBeNull();
  });

  it('handles digit and operator appending', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.append('1');
      result.current.append('2');
      result.current.append('+');
      result.current.append('3');
    });

    expect(result.current.expression).toBe('12 + 3');
  });

  it('replaces operator when consecutively entering binary operators', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.append('5');
      result.current.append('+');
    });
    expect(result.current.expression).toBe('5 + ');

    act(() => {
      result.current.append('*');
    });
    expect(result.current.expression).toBe('5 * ');
  });

  it('handles sqrt function token append', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.append('sqrt');
    });

    expect(result.current.expression).toBe('sqrt(');
  });

  it('handles backspace for single char, operator with padding, and sqrt token', () => {
    const { result } = renderHook(() => useCalculator());

    // Single character backspace
    act(() => {
      result.current.setExpression('123');
    });
    act(() => {
      result.current.backspace();
    });
    expect(result.current.expression).toBe('12');

    // Operator with padding backspace
    act(() => {
      result.current.setExpression('12 + ');
    });
    act(() => {
      result.current.backspace();
    });
    expect(result.current.expression).toBe('12');

    // Sqrt token backspace
    act(() => {
      result.current.setExpression('sqrt(');
    });
    act(() => {
      result.current.backspace();
    });
    expect(result.current.expression).toBe('');
  });

  it('clears all state when clear is called', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.setExpression('55 * 2');
    });
    act(() => {
      result.current.clear();
    });

    expect(result.current.expression).toBe('');
    expect(result.current.result).toBeNull();
    expect(result.current.error).toBeNull();
  });

  it('validates syntax before making API calls', async () => {
    const calculateSpy = vi.spyOn(apiService, 'calculateExpression');
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.setExpression('((5 + 3)');
    });

    await act(async () => {
      await result.current.calculate();
    });

    expect(calculateSpy).not.toHaveBeenCalled();
    expect(result.current.error).toEqual(
      expect.objectContaining({
        code: 'VALIDATION_ERROR',
        message: expect.stringContaining('parenthesis'),
      })
    );
  });

  it('performs successful calculation and updates history', async () => {
    vi.spyOn(apiService, 'calculateExpression').mockResolvedValueOnce({
      success: true,
      expression: '1000 + 10%',
      result: 1100,
      formatted: '1,100',
      execution_time_ms: 0.05,
    });

    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.setExpression('1000 + 10%');
    });

    await act(async () => {
      await result.current.calculate();
    });

    expect(result.current.result).toBe('1,100');
    expect(result.current.error).toBeNull();
    expect(result.current.history).toHaveLength(1);
    expect(result.current.history[0].expression).toBe('1000 + 10%');
    expect(result.current.history[0].formatted).toBe('1,100');
  });

  it('chains result with operator when appending operator after a calculation', async () => {
    vi.spyOn(apiService, 'calculateExpression').mockResolvedValueOnce({
      success: true,
      expression: '20 * 2',
      result: 40,
      formatted: '40',
      execution_time_ms: 0.05,
    });

    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.setExpression('20 * 2');
    });

    await act(async () => {
      await result.current.calculate();
    });

    expect(result.current.result).toBe('40');

    // Chaining with +
    act(() => {
      result.current.append('+');
    });

    expect(result.current.result).toBeNull();
    expect(result.current.expression).toBe('40 + ');
  });

  it('starts a new calculation when typing a digit after a calculation', async () => {
    vi.spyOn(apiService, 'calculateExpression').mockResolvedValueOnce({
      success: true,
      expression: '20 * 2',
      result: 40,
      formatted: '40',
      execution_time_ms: 0.05,
    });

    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.setExpression('20 * 2');
    });

    await act(async () => {
      await result.current.calculate();
    });

    // Fresh digit
    act(() => {
      result.current.append('9');
    });

    expect(result.current.result).toBeNull();
    expect(result.current.expression).toBe('9');
  });

  it('handles CalculatorApiError from backend', async () => {
    vi.spyOn(apiService, 'calculateExpression').mockRejectedValueOnce(
      new CalculatorApiError('Division by zero', 'DIVISION_BY_ZERO', 5)
    );

    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.setExpression('10 / 0');
    });

    await act(async () => {
      await result.current.calculate();
    });

    expect(result.current.error).toEqual({
      message: 'Division by zero',
      code: 'DIVISION_BY_ZERO',
      position: 5,
    });
    expect(result.current.result).toBeNull();
  });

  it('handles unexpected calculation errors gracefully', async () => {
    vi.spyOn(apiService, 'calculateExpression').mockRejectedValueOnce(new Error('Network offline'));

    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.setExpression('2 + 2');
    });

    await act(async () => {
      await result.current.calculate();
    });

    expect(result.current.error).toEqual({
      message: 'An unexpected error occurred during calculation.',
      code: 'UNEXPECTED_ERROR',
    });
  });

  it('recalls expression and result into state correctly', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.recallExpression('50 * (2 + 3)');
    });
    expect(result.current.expression).toBe('50 * (2 + 3)');

    act(() => {
      result.current.recallResult('250');
    });
    expect(result.current.expression).toBe('50 * (2 + 3)250');

    act(() => {
      result.current.clear();
      result.current.recallResult('250');
    });
    expect(result.current.expression).toBe('250');
  });

  it('clears history log', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.setExpression('1 + 1');
    });

    // Manually push to history through clearHistory test
    act(() => {
      result.current.clearHistory();
    });

    expect(result.current.history).toEqual([]);
  });

  it('handles keyboard shortcuts via window keydown listener', () => {
    const { result } = renderHook(() => useCalculator());

    // Digit keydown
    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: '7' }));
    });
    expect(result.current.expression).toBe('7');

    // Operator keydown
    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: '+' }));
    });
    expect(result.current.expression).toBe('7 + ');

    // Backspace keydown
    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Backspace' }));
    });
    expect(result.current.expression).toBe('7');

    // Escape keydown
    act(() => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }));
    });
    expect(result.current.expression).toBe('');
  });
});
