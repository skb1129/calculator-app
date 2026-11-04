import React from 'react';
import { Delete } from 'lucide-react';
import { KeyConfig } from '../types/calculator';

interface KeypadProps {
  onKeyPress: (val: string) => void;
  onClear: () => void;
  onBackspace: () => void;
  onCalculate: () => void;
  isLoading: boolean;
}

export const Keypad: React.FC<KeypadProps> = ({
  onKeyPress,
  onClear,
  onBackspace,
  onCalculate,
  isLoading,
}) => {
  const keys: KeyConfig[] = [
    // Row 1: Functional actions
    { label: 'AC', value: 'clear', type: 'action', ariaLabel: 'All clear' },
    { label: '⌫', value: 'backspace', type: 'action', ariaLabel: 'Backspace' },
    { label: '(', value: '(', type: 'operator', ariaLabel: 'Open parenthesis' },
    { label: ')', value: ')', type: 'operator', ariaLabel: 'Close parenthesis' },

    // Row 2: Advanced Math
    { label: '√x', value: 'sqrt', type: 'function', ariaLabel: 'Square root' },
    { label: 'xʸ', value: '^', type: 'function', ariaLabel: 'Exponentiation' },
    { label: '%', value: '%', type: 'function', ariaLabel: 'Percentage' },
    { label: '÷', value: '/', type: 'operator', ariaLabel: 'Division' },

    // Row 3: 7, 8, 9, Multiply
    { label: '7', value: '7', type: 'digit' },
    { label: '8', value: '8', type: 'digit' },
    { label: '9', value: '9', type: 'digit' },
    { label: '×', value: '*', type: 'operator', ariaLabel: 'Multiplication' },

    // Row 4: 4, 5, 6, Subtract
    { label: '4', value: '4', type: 'digit' },
    { label: '5', value: '5', type: 'digit' },
    { label: '6', value: '6', type: 'digit' },
    { label: '-', value: '-', type: 'operator', ariaLabel: 'Subtraction' },

    // Row 5: 1, 2, 3, Add
    { label: '1', value: '1', type: 'digit' },
    { label: '2', value: '2', type: 'digit' },
    { label: '3', value: '3', type: 'digit' },
    { label: '+', value: '+', type: 'operator', ariaLabel: 'Addition' },

    // Row 6: 0, Decimal, Equals
    { label: '0', value: '0', type: 'digit', span: 2 },
    { label: '.', value: '.', type: 'digit', ariaLabel: 'Decimal point' },
    { label: '=', value: 'equals', type: 'equals', ariaLabel: 'Calculate result' },
  ];

  const getKeyStyles = (type: KeyConfig['type']) => {
    switch (type) {
      case 'action':
        return 'bg-rose-950/40 text-rose-300 hover:bg-rose-900/60 active:bg-rose-900 border-rose-900/40';
      case 'function':
        return 'bg-purple-950/40 text-purple-300 hover:bg-purple-900/60 active:bg-purple-900 border-purple-900/40 font-semibold';
      case 'operator':
        return 'bg-amber-950/40 text-amber-300 hover:bg-amber-900/60 active:bg-amber-900 border-amber-900/40 text-xl font-semibold';
      case 'equals':
        return 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white hover:from-blue-500 hover:to-indigo-500 active:from-blue-700 active:to-indigo-700 font-bold text-2xl shadow-lg shadow-blue-500/25 border-blue-400/30';
      case 'digit':
      default:
        return 'bg-slate-800/80 text-slate-100 hover:bg-slate-700/80 active:bg-slate-700 border-slate-700/60 text-lg font-medium';
    }
  };

  const handleKeyClick = (key: KeyConfig) => {
    if (key.value === 'clear') {
      onClear();
    } else if (key.value === 'backspace') {
      onBackspace();
    } else if (key.value === 'equals') {
      onCalculate();
    } else {
      onKeyPress(key.value);
    }
  };

  return (
    <div className="grid grid-cols-4 gap-2 sm:gap-2.5">
      {keys.map((key) => {
        const spanClass = key.span === 2 ? 'col-span-2' : 'col-span-1';
        const isBackspace = key.value === 'backspace';

        return (
          <button
            key={key.label}
            onClick={() => handleKeyClick(key)}
            disabled={isLoading && key.value === 'equals'}
            aria-label={key.ariaLabel || key.label}
            className={`
              ${spanClass}
              ${getKeyStyles(key.type)}
              h-12 sm:h-14
              rounded-xl sm:rounded-2xl
              border
              flex items-center justify-center
              transition-all duration-150
              hover:scale-[1.02] active:scale-[0.98]
              focus:outline-none focus:ring-2 focus:ring-blue-500/50
              disabled:opacity-50 disabled:cursor-not-allowed
              select-none
            `}
          >
            {isBackspace ? (
              <Delete className="w-5 h-5" />
            ) : (
              <span>{key.label}</span>
            )}
          </button>
        );
      })}
    </div>
  );
};
