import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ErrorAlert } from '../src/components/ErrorAlert';

describe('ErrorAlert Component', () => {
  it('renders error message, code badge, and near-position text', () => {
    const error = {
      message: 'Unmatched opening parenthesis',
      code: 'UNMATCHED_PARENTHESIS',
      position: 3,
    };
    const onDismiss = vi.fn();

    render(<ErrorAlert error={error} onDismiss={onDismiss} />);

    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.getByText('Unmatched opening parenthesis')).toBeInTheDocument();
    expect(screen.getByText('UNMATCHED_PARENTHESIS')).toBeInTheDocument();
    expect(screen.getByText('Near position: 4')).toBeInTheDocument();
  });

  it('renders without position when position is undefined', () => {
    const error = {
      message: 'Expression exceeds maximum allowed length of 1000 characters',
      code: 'EXPRESSION_TOO_LONG',
    };
    const onDismiss = vi.fn();

    render(<ErrorAlert error={error} onDismiss={onDismiss} />);

    expect(screen.getByText('EXPRESSION_TOO_LONG')).toBeInTheDocument();
    expect(screen.queryByText(/Near position:/)).not.toBeInTheDocument();
  });

  it('calls onDismiss when close button is clicked', () => {
    const error = {
      message: 'Division by zero',
      code: 'DIVISION_BY_ZERO',
    };
    const onDismiss = vi.fn();

    render(<ErrorAlert error={error} onDismiss={onDismiss} />);

    const closeBtn = screen.getByRole('button', { name: 'Dismiss error' });
    fireEvent.click(closeBtn);

    expect(onDismiss).toHaveBeenCalledTimes(1);
  });
});
