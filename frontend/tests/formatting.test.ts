import { describe, it, expect } from 'vitest';
import { formatTimestamp, formatNumberWithCommas } from '../src/utils/formatting';

describe('Formatting Utils', () => {
  describe('formatTimestamp', () => {
    it('formats a valid ISO timestamp to locale time string', () => {
      const iso = '2026-11-04T12:30:45.000Z';
      const formatted = formatTimestamp(iso);
      expect(typeof formatted).toBe('string');
      expect(formatted.length).toBeGreaterThan(0);
    });

    it('falls back to returning the original string on invalid dates', () => {
      expect(formatTimestamp('not-a-date')).toBe('not-a-date');
    });
  });

  describe('formatNumberWithCommas', () => {
    it('adds commas to integers with 4 or more digits', () => {
      expect(formatNumberWithCommas('1000')).toBe('1,000');
      expect(formatNumberWithCommas('1000000')).toBe('1,000,000');
      expect(formatNumberWithCommas('123')).toBe('123');
    });

    it('formats floating point numbers without altering decimals', () => {
      expect(formatNumberWithCommas('12345.6789')).toBe('12,345.6789');
      expect(formatNumberWithCommas('0.123456')).toBe('0.123456');
    });

    it('preserves scientific notation strings unchanged', () => {
      expect(formatNumberWithCommas('1e+12')).toBe('1e+12');
      expect(formatNumberWithCommas('2.5E-6')).toBe('2.5E-6');
    });
  });
});
