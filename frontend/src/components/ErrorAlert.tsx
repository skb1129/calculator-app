import React from 'react';
import { AlertCircle, X } from 'lucide-react';

interface ErrorAlertProps {
  error: {
    message: string;
    code?: string;
    position?: number;
  };
  onDismiss: () => void;
}

export const ErrorAlert: React.FC<ErrorAlertProps> = ({ error, onDismiss }) => {
  return (
    <div
      role="alert"
      className="flex items-start gap-3 p-3.5 bg-rose-500/10 border border-rose-500/30 rounded-xl text-rose-300 text-sm animate-in fade-in slide-in-from-top-2 duration-200"
    >
      <AlertCircle className="w-5 h-5 text-rose-400 shrink-0 mt-0.5" />
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className="font-semibold text-rose-200">Calculation Error</span>
          {error.code && (
            <span className="text-xs px-2 py-0.5 bg-rose-500/20 text-rose-300 rounded font-mono">
              {error.code}
            </span>
          )}
        </div>
        <p className="mt-0.5 text-rose-300/90 break-words">{error.message}</p>
        {error.position !== undefined && (
          <p className="mt-1 text-xs text-rose-400/80">Near position: {error.position + 1}</p>
        )}
      </div>
      <button
        onClick={onDismiss}
        aria-label="Dismiss error"
        className="text-rose-400 hover:text-rose-200 transition-colors p-1 rounded-lg hover:bg-rose-500/20"
      >
        <X className="w-4 h-4" />
      </button>
    </div>
  );
};
