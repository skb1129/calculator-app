import { CalculateApiResponse, ApiError } from '../types/calculator';

const API_BASE = import.meta.env.VITE_API_BASE || '';

export class CalculatorApiError extends Error {
  code: string;
  position?: number;

  constructor(message: string, code: string = 'UNKNOWN_ERROR', position?: number) {
    super(message);
    this.name = 'CalculatorApiError';
    this.code = code;
    this.position = position;
  }
}

export async function calculateExpression(expression: string): Promise<CalculateApiResponse> {
  let response: Response;
  try {
    response = await fetch(`${API_BASE}/api/v1/calculate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ expression }),
    });
  } catch (err) {
    throw new CalculatorApiError(
      'Unable to connect to calculation backend service. Please check your connection.',
      'NETWORK_ERROR'
    );
  }

  const data = await response.json().catch(() => null);

  if (!response.ok) {
    if (data && typeof data === 'object' && 'error' in data) {
      const apiErr = data as ApiError;
      throw new CalculatorApiError(apiErr.error, apiErr.code, apiErr.position);
    }
    throw new CalculatorApiError(`Server error (${response.status})`, 'HTTP_ERROR');
  }

  return data as CalculateApiResponse;
}

export async function checkHealth(): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/health`);
    return res.ok;
  } catch {
    return false;
  }
}
