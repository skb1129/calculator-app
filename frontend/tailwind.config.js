/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        calc: {
          bg: '#0f172a',       // slate-900
          surface: '#1e293b',  // slate-800
          card: '#1e293b',
          border: '#334155',   // slate-700
          accent: '#3b82f6',   // blue-500
          accentHover: '#2563eb',
          operator: '#f59e0b', // amber-500
          operatorHover: '#d97706',
          function: '#8b5cf6', // purple-500
          functionHover: '#7c3aed',
          action: '#ef4444',   // red-500
          actionHover: '#dc2626',
        }
      },
      fontFamily: {
        mono: ['JetBrains Mono', 'Fira Code', 'Courier New', 'monospace'],
      }
    },
  },
  plugins: [],
}
