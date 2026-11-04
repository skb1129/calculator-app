import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { HistoryLog } from '../src/components/HistoryLog';
import { HistoryItem } from '../src/types/calculator';

describe('HistoryLog Component', () => {
  const mockHistory: HistoryItem[] = [
    {
      id: '1',
      expression: '1000 + 10%',
      result: 1100,
      formatted: '1,100',
      timestamp: '2026-11-04T12:00:00Z',
    },
    {
      id: '2',
      expression: 'sqrt(16) * 5',
      result: 20,
      formatted: '20',
      timestamp: '2026-11-04T12:05:00Z',
    },
  ];

  it('renders empty placeholder when history is empty', () => {
    render(
      <HistoryLog
        history={[]}
        onClearHistory={vi.fn()}
        onRecallExpression={vi.fn()}
        onRecallResult={vi.fn()}
      />
    );

    expect(screen.getByText('No calculations yet')).toBeInTheDocument();
    expect(screen.getByText('Calculations will be saved here in memory')).toBeInTheDocument();
    expect(screen.queryByTitle('Clear all calculation history')).not.toBeInTheDocument();
  });

  it('renders history items with expressions, formatted results, and badge count', () => {
    render(
      <HistoryLog
        history={mockHistory}
        onClearHistory={vi.fn()}
        onRecallExpression={vi.fn()}
        onRecallResult={vi.fn()}
      />
    );

    expect(screen.getByText('2')).toBeInTheDocument(); // count badge
    expect(screen.getByText('1000 + 10%')).toBeInTheDocument();
    expect(screen.getByText('= 1,100')).toBeInTheDocument();
    expect(screen.getByText('sqrt(16) * 5')).toBeInTheDocument();
    expect(screen.getByText('= 20')).toBeInTheDocument();
    expect(screen.getByTitle('Clear all calculation history')).toBeInTheDocument();
  });

  it('calls onClearHistory when Clear button is clicked', () => {
    const onClearHistory = vi.fn();
    render(
      <HistoryLog
        history={mockHistory}
        onClearHistory={onClearHistory}
        onRecallExpression={vi.fn()}
        onRecallResult={vi.fn()}
      />
    );

    const clearBtn = screen.getByTitle('Clear all calculation history');
    fireEvent.click(clearBtn);

    expect(onClearHistory).toHaveBeenCalledTimes(1);
  });

  it('recalls expression when clicking expression text or recall icon button', () => {
    const onRecallExpression = vi.fn();
    render(
      <HistoryLog
        history={mockHistory}
        onClearHistory={vi.fn()}
        onRecallExpression={onRecallExpression}
        onRecallResult={vi.fn()}
      />
    );

    // Click expression text
    const exprText = screen.getByText('1000 + 10%');
    fireEvent.click(exprText);
    expect(onRecallExpression).toHaveBeenCalledWith('1000 + 10%');

    // Click recall button icon
    const recallButtons = screen.getAllByTitle('Load expression back into editor');
    fireEvent.click(recallButtons[1]);
    expect(onRecallExpression).toHaveBeenCalledWith('sqrt(16) * 5');
  });

  it('recalls result when clicking result text or insert icon button', () => {
    const onRecallResult = vi.fn();
    render(
      <HistoryLog
        history={mockHistory}
        onClearHistory={vi.fn()}
        onRecallExpression={vi.fn()}
        onRecallResult={onRecallResult}
      />
    );

    // Click result text
    const resultText = screen.getByText('= 1,100');
    fireEvent.click(resultText);
    expect(onRecallResult).toHaveBeenCalledWith('1,100');

    // Click insert result icon button
    const insertButtons = screen.getAllByTitle('Insert result into active calculation');
    fireEvent.click(insertButtons[1]);
    expect(onRecallResult).toHaveBeenCalledWith('20');
  });
});
