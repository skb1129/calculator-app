import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Calculator } from '../src/components/Calculator';

describe('Calculator Component', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('renders display and keypad with all primary buttons', () => {
    render(<Calculator />);

    expect(screen.getByLabelText('Calculator expression input')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'All clear' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Calculate result' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Square root' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Exponentiation' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Percentage' })).toBeInTheDocument();
  });

  it('updates display when keypad buttons are clicked', () => {
    render(<Calculator />);
    const input = screen.getByLabelText('Calculator expression input') as HTMLInputElement;

    fireEvent.click(screen.getByRole('button', { name: '7' }));
    fireEvent.click(screen.getByRole('button', { name: 'Addition' }));
    fireEvent.click(screen.getByRole('button', { name: '8' }));

    expect(input.value).toBe('7 + 8');
  });

  it('performs calculation via API call and updates history', async () => {
    // Mock global fetch for API
    const mockApiResponse = {
      success: true,
      expression: '5 * 4',
      result: 20,
      formatted: '20',
      execution_time_ms: 0.12,
    };

    vi.spyOn(global, 'fetch').mockImplementationOnce(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockApiResponse),
      } as Response)
    );

    render(<Calculator />);
    const input = screen.getByLabelText('Calculator expression input') as HTMLInputElement;

    fireEvent.change(input, { target: { value: '5 * 4' } });
    fireEvent.click(screen.getByRole('button', { name: 'Calculate result' }));

    await waitFor(() => {
      expect(screen.getByTestId('display-result')).toHaveTextContent('= 20');
    });

    // Verify item added to history
    expect(screen.getAllByText('5 * 4').length).toBeGreaterThan(0);
  });

  it('displays error alert when backend returns calculation error', async () => {
    const mockErrorResponse = {
      success: false,
      error: 'Division by zero at position 4',
      code: 'DIVISION_BY_ZERO',
      position: 4,
    };

    vi.spyOn(global, 'fetch').mockImplementationOnce(() =>
      Promise.resolve({
        ok: false,
        status: 422,
        json: () => Promise.resolve(mockErrorResponse),
      } as Response)
    );

    render(<Calculator />);
    const input = screen.getByLabelText('Calculator expression input') as HTMLInputElement;

    fireEvent.change(input, { target: { value: '10 / 0' } });
    fireEvent.click(screen.getByRole('button', { name: 'Calculate result' }));

    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
      expect(screen.getByText('Division by zero at position 4')).toBeInTheDocument();
      expect(screen.getByText('DIVISION_BY_ZERO')).toBeInTheDocument();
    });
  });

  it('clears display and error when AC is clicked', async () => {
    render(<Calculator />);
    const input = screen.getByLabelText('Calculator expression input') as HTMLInputElement;

    fireEvent.change(input, { target: { value: '123 + 456' } });
    expect(input.value).toBe('123 + 456');

    fireEvent.click(screen.getByRole('button', { name: 'All clear' }));
    expect(input.value).toBe('');
  });
});
