export interface HistoryItem {
  id: string;
  expression: string;
  result: number;
  formatted: string;
  timestamp: string;
}

export interface CalculateApiResponse {
  success: boolean;
  expression: string;
  result: number;
  formatted: string;
  tokens?: string[];
  execution_time_ms: number;
}

export interface ApiError {
  success: false;
  error: string;
  code: string;
  position?: number;
}

export type KeyType = 'digit' | 'operator' | 'function' | 'action' | 'equals';

export interface KeyConfig {
  label: string;
  value: string;
  type: KeyType;
  ariaLabel?: string;
  span?: number;
}
