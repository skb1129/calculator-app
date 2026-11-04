import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { Display } from '../src/components/Display';

describe('Display Component', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('renders expression and default ready status when no result', () => {
    render(
      <Display
        expression="10 + 20"
        result={null}
        isLoading={false}
        onExpressionChange={vi.fn()}
        onSubmit={vi.fn()}
      />
    );

    const input = screen.getByLabelText('Calculator expression input') as HTMLInputElement;
    expect(input.value).toBe('10 + 20');
    expect(screen.getByText('Ready')).toBeInTheDocument();
    expect(screen.queryByText('Copy')).not.toBeInTheDocument();
  });

  it('calls onExpressionChange when typing in the input field', () => {
    const onExpressionChange = vi.fn();
    render(
      <Display
        expression=""
        result={null}
        isLoading={false}
        onExpressionChange={onExpressionChange}
        onSubmit={vi.fn()}
      />
    );

    const input = screen.getByLabelText('Calculator expression input');
    fireEvent.change(input, { target: { value: '42 * 2' } });

    expect(onExpressionChange).toHaveBeenCalledWith('42 * 2');
  });

  it('triggers onSubmit when Enter key is pressed in the input', () => {
    const onSubmit = vi.fn();
    render(
      <Display
        expression="100 / 5"
        result={null}
        isLoading={false}
        onExpressionChange={vi.fn()}
        onSubmit={onSubmit}
      />
    );

    const input = screen.getByLabelText('Calculator expression input');
    fireEvent.keyDown(input, { key: 'Enter', code: 'Enter' });

    expect(onSubmit).toHaveBeenCalledTimes(1);
  });

  it('shows animated loading spinner and text while calculating', () => {
    render(
      <Display
        expression="sqrt(144)"
        result={null}
        isLoading={true}
        onExpressionChange={vi.fn()}
        onSubmit={vi.fn()}
      />
    );

    expect(screen.getByText('Calculating server-side...')).toBeInTheDocument();
  });

  it('renders calculated result when result is provided', () => {
    render(
      <Display
        expression="25 * 4"
        result="100"
        isLoading={false}
        onExpressionChange={vi.fn()}
        onSubmit={vi.fn()}
      />
    );

    const resultSpan = screen.getByTestId('display-result');
    expect(resultSpan).toHaveTextContent('= 100');
  });

  it('copies result to clipboard and shows temporary Copied feedback', async () => {
    const writeTextMock = vi.fn().mockResolvedValue(undefined);
    Object.assign(navigator, {
      clipboard: {
        writeText: writeTextMock,
      },
    });

    vi.useFakeTimers();

    render(
      <Display
        expression="50 + 50"
        result="100"
        isLoading={false}
        onExpressionChange={vi.fn()}
        onSubmit={vi.fn()}
      />
    );

    const copyBtn = screen.getByRole('button', { name: /copy/i });
    expect(copyBtn).toBeInTheDocument();

    fireEvent.click(copyBtn);

    expect(writeTextMock).toHaveBeenCalledWith('100');
    expect(screen.getByText('Copied!')).toBeInTheDocument();

    // Fast-forward 1500ms
    act(() => {
      vi.advanceTimersByTime(1500);
    });

    expect(screen.queryByText('Copied!')).not.toBeInTheDocument();
    expect(screen.getByText('Copy')).toBeInTheDocument();

    vi.useRealTimers();
  });
});
