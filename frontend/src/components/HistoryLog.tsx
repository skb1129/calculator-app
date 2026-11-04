import React from 'react';
import { History, Trash2, ArrowUpRight, PlusCircle } from 'lucide-react';
import { HistoryItem } from '../types/calculator';
import { formatTimestamp } from '../utils/formatting';

interface HistoryLogProps {
  history: HistoryItem[];
  onClearHistory: () => void;
  onRecallExpression: (expr: string) => void;
  onRecallResult: (res: string) => void;
}

export const HistoryLog: React.FC<HistoryLogProps> = ({
  history,
  onClearHistory,
  onRecallExpression,
  onRecallResult,
}) => {
  return (
    <div className="flex flex-col h-full bg-slate-900/60 rounded-2xl border border-slate-800 p-4 sm:p-5">
      {/* Header */}
      <div className="flex items-center justify-between pb-3 border-b border-slate-800">
        <div className="flex items-center gap-2">
          <History className="w-4 h-4 text-blue-400" />
          <h2 className="text-sm font-semibold uppercase tracking-wider text-slate-300">
            Calculation History
          </h2>
          <span className="text-xs px-2 py-0.5 bg-slate-800 text-slate-400 rounded-full font-mono">
            {history.length}
          </span>
        </div>

        {history.length > 0 && (
          <button
            onClick={onClearHistory}
            title="Clear all calculation history"
            className="flex items-center gap-1.5 text-xs text-rose-400 hover:text-rose-300 transition-colors p-1.5 rounded-lg hover:bg-rose-500/10 active:scale-95"
          >
            <Trash2 className="w-3.5 h-3.5" />
            <span>Clear</span>
          </button>
        )}
      </div>

      {/* History item list */}
      <div className="flex-1 overflow-y-auto mt-3 pr-1 space-y-2.5 max-h-[609px] scrollbar-thin scrollbar-thumb-slate-700">
        {history.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-48 text-center text-slate-500">
            <History className="w-8 h-8 mb-2 stroke-1 text-slate-600" />
            <p className="text-sm">No calculations yet</p>
            <p className="text-xs text-slate-600 mt-1">Calculations will be saved here in memory</p>
          </div>
        ) : (
          history.map((item) => (
            <div
              key={item.id}
              className="group p-3 rounded-xl bg-slate-800/40 hover:bg-slate-800/80 border border-slate-800/60 hover:border-slate-700 transition-all text-left"
            >
              <div className="flex items-center justify-between text-xs text-slate-500">
                <span className="font-mono">{formatTimestamp(item.timestamp)}</span>
                <div className="flex items-center gap-1 opacity-80 group-hover:opacity-100 transition-opacity">
                  <button
                    onClick={() => onRecallExpression(item.expression)}
                    title="Load expression back into editor"
                    className="p-1 text-slate-400 hover:text-blue-400 rounded hover:bg-slate-700/50 transition-colors"
                  >
                    <ArrowUpRight className="w-3.5 h-3.5" />
                  </button>
                  <button
                    onClick={() => onRecallResult(item.formatted)}
                    title="Insert result into active calculation"
                    className="p-1 text-slate-400 hover:text-emerald-400 rounded hover:bg-slate-700/50 transition-colors"
                  >
                    <PlusCircle className="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>

              <div
                onClick={() => onRecallExpression(item.expression)}
                className="mt-1 font-mono text-sm text-slate-300 break-words cursor-pointer hover:text-blue-300 transition-colors"
              >
                {item.expression}
              </div>

              <div
                onClick={() => onRecallResult(item.formatted)}
                className="mt-0.5 font-mono text-base font-bold text-emerald-400 text-right cursor-pointer hover:underline"
              >
                = {item.formatted}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};
