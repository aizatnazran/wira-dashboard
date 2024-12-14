/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'ac-gold': '#C6A875',
        'ac-dark': '#1A1A1A',
        'ac-gray': '#2A2A2A',
        'ac-light': '#F5F5F5',
      },
      fontFamily: {
        'cinzel': ['Cinzel', 'serif'],
        'montserrat': ['Montserrat', 'sans-serif'],
      },
      backgroundImage: {
        'hero-pattern': "linear-gradient(to bottom, rgba(0, 0, 0, 0.7), rgba(0, 0, 0, 0.5))",
      },
    },
  },
  plugins: [],
}
