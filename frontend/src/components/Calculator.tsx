import React, { useState } from 'react';
import { History, Sparkles } from 'lucide-react';
import { useCalculator } from '../hooks/useCalculator';
import { Display } from './Display';
import { Keypad } from './Keypad';
import { HistoryLog } from './HistoryLog';
import { ErrorAlert } from './ErrorAlert';

export const Calculator: React.FC = () => {
  const {
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
  } = useCalculator();

  const [isMobileHistoryOpen, setIsMobileHistoryOpen] = useState(false);

  const sampleExpressions = [
    '(15 + 5 * 2^3) / sqrt(16) - 10%',
    '100 * 15% + 50',
    '2 ^ 3 ^ 2',
    'sqrt(144) + √(49)',
  ];

  return (
    <div className="w-full max-w-5xl mx-auto">
      {/* Top action / mobile toggle bar */}
      <div className="flex items-center justify-between mb-4 lg:hidden">
        <div className="flex items-center gap-2">
          <span className="w-2.5 h-2.5 rounded-full bg-emerald-400 animate-pulse"></span>
          <span className="text-xs text-slate-400 font-mono">Backend API Connected</span>
        </div>
        <button
          onClick={() => setIsMobileHistoryOpen(!isMobileHistoryOpen)}
          className="flex items-center gap-1.5 text-xs text-slate-300 bg-slate-800 hover:bg-slate-700 px-3 py-1.5 rounded-lg border border-slate-700 transition-colors"
        >
          <History className="w-3.5 h-3.5 text-blue-400" />
          <span>History ({history.length})</span>
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        {/* Calculator Main Column (7 cols on desktop) */}
        <div className="lg:col-span-7 bg-slate-900/80 rounded-3xl p-5 sm:p-6 border border-slate-800 shadow-2xl backdrop-blur-sm space-y-4">
          {/* Error Alert */}
          {error && (
            <ErrorAlert
              error={error}
              onDismiss={() => {
                /* Hook resets error on next keypress */
              }}
            />
          )}

          {/* Calculator Display */}
          <Display
            expression={expression}
            result={result}
            isLoading={isLoading}
            onExpressionChange={setExpression}
            onSubmit={calculate}
          />

          {/* Keypad */}
          <Keypad
            onKeyPress={append}
            onClear={clear}
            onBackspace={backspace}
            onCalculate={calculate}
            isLoading={isLoading}
          />

          {/* Preset sample formulas for rapid testing */}
          <div className="pt-2 border-t border-slate-800/80">
            <div className="flex items-center gap-1.5 text-xs text-slate-400 mb-2">
              <Sparkles className="w-3 h-3 text-amber-400" />
              <span>Sample Expressions:</span>
            </div>
            <div className="flex flex-wrap gap-1.5">
              {sampleExpressions.map((sample) => (
                <button
                  key={sample}
                  onClick={() => recallExpression(sample)}
                  className="text-xs font-mono bg-slate-800/60 hover:bg-slate-700 text-slate-300 hover:text-white px-2.5 py-1 rounded-md border border-slate-700/50 transition-all hover:scale-[1.02]"
                >
                  {sample}
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* History Column (5 cols on desktop, responsive drawer on mobile) */}
        <div className="hidden lg:block lg:col-span-5 h-full">
          <HistoryLog
            history={history}
            onClearHistory={clearHistory}
            onRecallExpression={recallExpression}
            onRecallResult={recallResult}
          />
        </div>
      </div>

      {/* Mobile History Bottom Drawer */}
      {isMobileHistoryOpen && (
        <div className="fixed inset-0 z-50 lg:hidden flex flex-col justify-end bg-black/60 backdrop-blur-sm">
          <div
            className="fixed inset-0"
            onClick={() => setIsMobileHistoryOpen(false)}
          ></div>
          <div className="relative z-10 bg-slate-900 rounded-t-3xl border-t border-slate-700 p-5 max-h-[80vh] flex flex-col animate-in slide-in-from-bottom duration-200 shadow-2xl">
            <div className="w-12 h-1 bg-slate-700 rounded-full mx-auto mb-4"></div>
            <HistoryLog
              history={history}
              onClearHistory={clearHistory}
              onRecallExpression={(expr) => {
                recallExpression(expr);
                setIsMobileHistoryOpen(false);
              }}
              onRecallResult={(res) => {
                recallResult(res);
                setIsMobileHistoryOpen(false);
              }}
            />
          </div>
        </div>
      )}
    </div>
  );
};
