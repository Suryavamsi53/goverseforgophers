/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./ui/templates/**/*.html",
    "./ui/assets/**/*.js",
  ],
  theme: {
    extend: {
      colors: {
        go: {
          cyan: '#00ADD8',
          blue: '#007D9C',
          light: '#E0F6FB',
          dark: '#005367',
        },
        dark: {
          900: '#0F172A', // Slate 900
          800: '#1E293B', // Slate 800
          700: '#334155', // Slate 700
          border: '#334155',
        }
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
      },
    },
  },
  plugins: [],
}
