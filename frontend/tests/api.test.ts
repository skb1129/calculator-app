import { describe, it, expect, vi, beforeEach } from 'vitest';
import { calculateExpression, checkHealth, CalculatorApiError } from '../src/services/api';

describe('API Service', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  describe('calculateExpression', () => {
    it('returns CalculateApiResponse on successful HTTP 200', async () => {
      const mockResponse = {
        success: true,
        expression: '1000 + 10%',
        result: 1100,
        formatted: '1100',
        execution_time_ms: 0.05,
      };

      vi.spyOn(global, 'fetch').mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      } as Response);

      const res = await calculateExpression('1000 + 10%');
      expect(res.success).toBe(true);
      expect(res.result).toBe(1100);
      expect(res.formatted).toBe('1100');
    });

    it('throws CalculatorApiError with code and position on HTTP 400/422 error', async () => {
      const mockErrorResponse = {
        success: false,
        error: 'Division by zero at position 5',
        code: 'DIVISION_BY_ZERO',
        position: 5,
      };

      vi.spyOn(global, 'fetch').mockResolvedValueOnce({
        ok: false,
        status: 422,
        json: () => Promise.resolve(mockErrorResponse),
      } as Response);

      await expect(calculateExpression('10 / 0')).rejects.toThrow('Division by zero at position 5');
    });

    it('throws CalculatorApiError with NETWORK_ERROR on connection failure', async () => {
      vi.spyOn(global, 'fetch').mockRejectedValueOnce(new Error('Failed to fetch'));

      await expect(calculateExpression('5 + 5')).rejects.toThrow(CalculatorApiError);
    });
  });

  describe('checkHealth', () => {
    it('returns true when health endpoint responds with 200 OK', async () => {
      vi.spyOn(global, 'fetch').mockResolvedValueOnce({
        ok: true,
      } as Response);

      const isHealthy = await checkHealth();
      expect(isHealthy).toBe(true);
    });

    it('returns false when health endpoint fails or throws', async () => {
      vi.spyOn(global, 'fetch').mockRejectedValueOnce(new Error('Network error'));

      const isHealthy = await checkHealth();
      expect(isHealthy).toBe(false);
    });
  });
});
