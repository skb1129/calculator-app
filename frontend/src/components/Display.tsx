import React, { useRef, useEffect, useState } from 'react';
import { Copy, Check, Loader2 } from 'lucide-react';

interface DisplayProps {
  expression: string;
  result: string | null;
  isLoading: boolean;
  onExpressionChange: (val: string) => void;
  onSubmit: () => void;
}

export const Display: React.FC<DisplayProps> = ({
  expression,
  result,
  isLoading,
  onExpressionChange,
  onSubmit,
}) => {
  const [copied, setCopied] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    // Scroll input to end when expression changes
    if (inputRef.current) {
      inputRef.current.scrollLeft = inputRef.current.scrollWidth;
    }
  }, [expression]);

  const handleCopy = () => {
    if (result) {
      navigator.clipboard.writeText(result);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    }
  };

  return (
    <div className="bg-slate-900/90 rounded-2xl p-4 sm:p-5 border border-slate-800 shadow-inner flex flex-col justify-between min-h-[130px] sm:min-h-[145px] transition-all">
      {/* Expression input line */}
      <div className="relative flex items-center">
        <input
          ref={inputRef}
          type="text"
          value={expression}
          onChange={(e) => onExpressionChange(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              e.preventDefault();
              onSubmit();
            }
          }}
          placeholder="0"
          aria-label="Calculator expression input"
          className="w-full bg-transparent text-slate-300 font-mono text-lg sm:text-xl md:text-2xl outline-none placeholder:text-slate-600 tracking-wider text-right"
        />
      </div>

      {/* Result line */}
      <div className="flex items-end justify-between mt-3 pt-2 border-t border-slate-800/60">
        <div className="flex items-center gap-2">
          {isLoading ? (
            <div className="flex items-center gap-2 text-xs text-blue-400 font-medium animate-pulse">
              <Loader2 className="w-3.5 h-3.5 animate-spin" />
              <span>Calculating server-side...</span>
            </div>
          ) : result !== null ? (
            <button
              onClick={handleCopy}
              title="Copy result to clipboard"
              className="flex items-center gap-1.5 text-xs text-slate-400 hover:text-slate-200 transition-colors py-1 px-2 rounded-md hover:bg-slate-800/80 active:scale-95"
            >
              {copied ? (
                <>
                  <Check className="w-3.5 h-3.5 text-emerald-400" />
                  <span className="text-emerald-400 font-medium">Copied!</span>
                </>
              ) : (
                <>
                  <Copy className="w-3.5 h-3.5" />
                  <span>Copy</span>
                </>
              )}
            </button>
          ) : (
            <span className="text-xs text-slate-500 font-mono">Ready</span>
          )}
        </div>

        <div className="text-right">
          <span
            data-testid="display-result"
            className={`font-mono font-bold tracking-tight select-all transition-all ${
              result !== null
                ? 'text-3xl sm:text-4xl text-emerald-400'
                : 'text-2xl text-slate-600'
            }`}
          >
            {result !== null ? `= ${result}` : ''}
          </span>
        </div>
      </div>
    </div>
  );
};
