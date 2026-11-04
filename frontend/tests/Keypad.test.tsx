import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { Keypad } from '../src/components/Keypad';

describe('Keypad Component', () => {
  it('renders all 23 keypad buttons correctly', () => {
    render(
      <Keypad
        onKeyPress={vi.fn()}
        onClear={vi.fn()}
        onBackspace={vi.fn()}
        onCalculate={vi.fn()}
        isLoading={false}
      />
    );

    const buttons = screen.getAllByRole('button');
    expect(buttons).toHaveLength(23);
  });

  it('triggers onKeyPress with character value when number or operator is clicked', () => {
    const onKeyPress = vi.fn();
    render(
      <Keypad
        onKeyPress={onKeyPress}
        onClear={vi.fn()}
        onBackspace={vi.fn()}
        onCalculate={vi.fn()}
        isLoading={false}
      />
    );

    fireEvent.click(screen.getByRole('button', { name: '9' }));
    expect(onKeyPress).toHaveBeenCalledWith('9');

    fireEvent.click(screen.getByRole('button', { name: 'Addition' }));
    expect(onKeyPress).toHaveBeenCalledWith('+');

    fireEvent.click(screen.getByRole('button', { name: 'Square root' }));
    expect(onKeyPress).toHaveBeenCalledWith('sqrt');

    fireEvent.click(screen.getByRole('button', { name: 'Exponentiation' }));
    expect(onKeyPress).toHaveBeenCalledWith('^');

    fireEvent.click(screen.getByRole('button', { name: 'Percentage' }));
    expect(onKeyPress).toHaveBeenCalledWith('%');
  });

  it('triggers onClear when AC button is clicked', () => {
    const onClear = vi.fn();
    render(
      <Keypad
        onKeyPress={vi.fn()}
        onClear={onClear}
        onBackspace={vi.fn()}
        onCalculate={vi.fn()}
        isLoading={false}
      />
    );

    fireEvent.click(screen.getByRole('button', { name: 'All clear' }));
    expect(onClear).toHaveBeenCalledTimes(1);
  });

  it('triggers onBackspace when ⌫ button is clicked', () => {
    const onBackspace = vi.fn();
    render(
      <Keypad
        onKeyPress={vi.fn()}
        onClear={vi.fn()}
        onBackspace={onBackspace}
        onCalculate={vi.fn()}
        isLoading={false}
      />
    );

    fireEvent.click(screen.getByRole('button', { name: 'Backspace' }));
    expect(onBackspace).toHaveBeenCalledTimes(1);
  });

  it('triggers onCalculate when = button is clicked', () => {
    const onCalculate = vi.fn();
    render(
      <Keypad
        onKeyPress={vi.fn()}
        onClear={vi.fn()}
        onBackspace={vi.fn()}
        onCalculate={onCalculate}
        isLoading={false}
      />
    );

    fireEvent.click(screen.getByRole('button', { name: 'Calculate result' }));
    expect(onCalculate).toHaveBeenCalledTimes(1);
  });

  it('disables the = button when isLoading is true', () => {
    render(
      <Keypad
        onKeyPress={vi.fn()}
        onClear={vi.fn()}
        onBackspace={vi.fn()}
        onCalculate={vi.fn()}
        isLoading={true}
      />
    );

    const equalsBtn = screen.getByRole('button', { name: 'Calculate result' });
    expect(equalsBtn).toBeDisabled();
  });
});
