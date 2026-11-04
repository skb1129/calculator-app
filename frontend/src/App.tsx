import { useState, useEffect } from 'react';
import { Calculator as CalcIcon, Cpu, Code2 } from 'lucide-react';
import { Calculator } from './components/Calculator';
import { checkHealth } from './services/api';

export function App() {
  const [isBackendHealthy, setIsBackendHealthy] = useState<boolean | null>(null);

  useEffect(() => {
    checkHealth().then(setIsBackendHealthy);
    const interval = setInterval(() => {
      checkHealth().then(setIsBackendHealthy);
    }, 15000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-950 via-slate-900 to-slate-950 text-slate-100 flex flex-col justify-between py-6 px-4 sm:px-6 lg:px-8">
      {/* Header */}
      <header className="max-w-5xl mx-auto w-full mb-6">
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4 pb-5 border-b border-slate-800/80">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center shadow-lg shadow-blue-500/20">
              <CalcIcon className="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 className="text-xl sm:text-2xl font-bold tracking-tight text-white flex items-center gap-2">
                OmniCalc
                <span className="text-xs font-mono font-medium px-2 py-0.5 rounded-full bg-blue-500/10 text-blue-400 border border-blue-500/20">
                  Full-Stack
                </span>
              </h1>
              <p className="text-xs text-slate-400">
                Multi-operation expression evaluator powered by Go Gin backend
              </p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            {/* Health Status Indicator */}
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-slate-900 border border-slate-800 text-xs">
              <Cpu className="w-3.5 h-3.5 text-slate-400" />
              <span className="text-slate-400">Go Backend:</span>
              {isBackendHealthy === null ? (
                <span className="text-amber-400 flex items-center gap-1">
                  <span className="w-1.5 h-1.5 rounded-full bg-amber-400 animate-ping"></span>
                  Checking
                </span>
              ) : isBackendHealthy ? (
                <span className="text-emerald-400 font-medium flex items-center gap-1.5">
                  <span className="w-2 h-2 rounded-full bg-emerald-400"></span>
                  Connected
                </span>
              ) : (
                <span className="text-rose-400 font-medium flex items-center gap-1.5">
                  <span className="w-2 h-2 rounded-full bg-rose-400"></span>
                  Disconnected
                </span>
              )}
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 flex flex-col justify-center max-w-5xl mx-auto w-full">
        <Calculator />
      </main>

      {/* Footer */}
      <footer className="max-w-5xl mx-auto w-full mt-8 pt-5 border-t border-slate-800/80 text-xs text-slate-500 flex flex-col sm:flex-row items-center justify-between gap-3">
        <div className="flex flex-wrap items-center justify-center gap-2">
          <span className="text-slate-400">Supported Operations:</span>
          <span className="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300 font-mono">Addition (+)</span>
          <span className="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300 font-mono">Subtraction (-)</span>
          <span className="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300 font-mono">Multiplication (*)</span>
          <span className="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300 font-mono">Division (/)</span>
          <span className="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300 font-mono">Exponentiation (^)</span>
          <span className="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300 font-mono">Square Root (sqrt)</span>
          <span className="px-1.5 py-0.5 rounded bg-slate-800 text-slate-300 font-mono">Percentage (%)</span>
        </div>

        <div className="flex items-center gap-2">
          <span>Clean &amp; Idiomatic Architecture</span>
          <span>•</span>
          <div className="flex items-center gap-1 text-slate-400">
            <Code2 className="w-3.5 h-3.5" />
            <span>React + Go Gin</span>
          </div>
        </div>
      </footer>
    </div>
  );
}

export default App;
