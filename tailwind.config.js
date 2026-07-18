/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./ui/templates/**/*.html",
    "./ui/assets/**/*.js",
  ],
  darkMode: 'class', // Force dark mode by default if desired
  theme: {
    extend: {
      colors: {
        app: '#000000',
        panel: '#000000',
        elevated: '#000000',
        glow: '#00D084',
        glowHover: '#00E691',
        glowMuted: 'rgba(0, 208, 132, 0.15)',
        primary: '#EDEDED',
        secondary: '#A1A1AA',
        tertiary: '#52525B',
        success: '#238636',
        error: '#F85149',
        warning: '#D29922',
        info: '#2F81F7',
      },
      fontFamily: {
        sans: ['Inter', 'SF Pro Display', 'sans-serif'],
        mono: ['JetBrains Mono', 'Geist Mono', 'monospace'],
      },
      boxShadow: {
        glow: '0 0 15px rgba(0, 208, 132, 0.2)',
        modal: '0 25px 50px -12px rgba(0, 0, 0, 0.5)',
      },
      animation: {
        'pulse-glow': 'pulse-glow 2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
      keyframes: {
        'pulse-glow': {
          '0%, 100%': { opacity: 1, boxShadow: '0 0 15px rgba(0, 208, 132, 0.2)' },
          '50%': { opacity: .7, boxShadow: '0 0 5px rgba(0, 208, 132, 0.1)' },
        }
      }
    },
  },
  plugins: [],
}
