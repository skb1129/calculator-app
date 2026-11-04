import { useState, useCallback, useEffect } from 'react';
import { HistoryItem } from '../types/calculator';
import { calculateExpression, CalculatorApiError } from '../services/api';
import { validateExpression } from '../utils/validation';

export function useCalculator() {
  const [expression, setExpression] = useState<string>('');
  const [result, setResult] = useState<string | null>(null);
  const [history, setHistory] = useState<HistoryItem[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<{ message: string; code?: string; position?: number } | null>(null);

  const append = useCallback((val: string) => {
    setError(null);
    setExpression((prev) => {
      // If we previously had a result and the user enters a binary operator, chain it with previous result!
      if (result !== null && ['+', '-', '*', '/', '^', '%'].includes(val)) {
        const chained = `${result} ${val} `;
        setResult(null);
        return chained;
      }

      // If user starts typing a fresh number after a result, start fresh
      if (result !== null && /[0-9(]/.test(val)) {
        setResult(null);
        return val;
      }

      // Handle function syntax like "sqrt("
      if (val === 'sqrt') {
        return prev + 'sqrt(';
      }

      // Add clean spacing for standard binary operators
      if (['+', '-', '*', '/'].includes(val)) {
        // If already ending with space, don't duplicate
        if (prev.endsWith(' ')) {
          return prev.slice(0, -2) + val + ' ';
        }
        return prev + ` ${val} `;
      }

      return prev + val;
    });
  }, [result]);

  const clear = useCallback(() => {
    setExpression('');
    setResult(null);
    setError(null);
  }, []);

  const backspace = useCallback(() => {
    setError(null);
    setExpression((prev) => {
      if (!prev) return '';
      // If ending with "sqrt(", remove the entire token
      if (prev.endsWith('sqrt(')) {
        return prev.slice(0, -5);
      }
      // If ending with operator and trailing space " + ", remove " + "
      if (/\s[+\-*/]\s$/.test(prev)) {
        return prev.slice(0, -3);
      }
      return prev.slice(0, -1);
    });
  }, []);

  const calculate = useCallback(async () => {
    if (!expression.trim()) return;

    setError(null);

    // Client-side quick syntax validation
    const validation = validateExpression(expression);
    if (!validation.isValid) {
      setError({
        message: validation.error || 'Invalid expression syntax',
        code: 'VALIDATION_ERROR',
        position: validation.position,
      });
      return;
    }

    setIsLoading(true);
    try {
      const data = await calculateExpression(expression);
      setResult(data.formatted);

      // Prepend to in-memory history log
      const historyEntry: HistoryItem = {
        id: (Date.now() + Math.random()).toString(),
        expression: expression.trim(),
        result: data.result,
        formatted: data.formatted,
        timestamp: new Date().toISOString(),
      };
      setHistory((prev) => [historyEntry, ...prev]);
    } catch (err) {
      if (err instanceof CalculatorApiError) {
        setError({
          message: err.message,
          code: err.code,
          position: err.position,
        });
      } else {
        setError({
          message: 'An unexpected error occurred during calculation.',
          code: 'UNEXPECTED_ERROR',
        });
      }
    } finally {
      setIsLoading(false);
    }
  }, [expression]);

  const clearHistory = useCallback(() => {
    setHistory([]);
  }, []);

  const recallExpression = useCallback((expr: string) => {
    setExpression(expr);
    setResult(null);
    setError(null);
  }, []);

  const recallResult = useCallback((val: string) => {
    setExpression((prev) => {
      if (!prev || result !== null) {
        setResult(null);
        return val;
      }
      return prev + val;
    });
    setError(null);
  }, [result]);

  // Global physical keyboard support
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Don't intercept if user is focusing an external input or selecting text
      const target = e.target as HTMLElement;
      if (target.tagName === 'INPUT' && target.getAttribute('type') === 'text') {
        return;
      }

      if (e.key === 'Enter' || e.key === '=') {
        e.preventDefault();
        calculate();
      } else if (e.key === 'Backspace') {
        e.preventDefault();
        backspace();
      } else if (e.key === 'Escape') {
        e.preventDefault();
        clear();
      } else if (/[0-9]/.test(e.key)) {
        e.preventDefault();
        append(e.key);
      } else if (['+', '-', '*', '/', '^', '%', '(', ')', '.'].includes(e.key)) {
        e.preventDefault();
        append(e.key);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [append, backspace, calculate, clear]);

  return {
    expression,
    setExpression,
    result,
    history,
    isLoading,
    error,
    append,
    clear,
    backspace,
    calculate,
    clearHistory,
    recallExpression,
    recallResult,
  };
}
