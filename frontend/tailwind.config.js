/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: "#646cff",
        dark: "#0f0f0f",
        "dark-card": "#1a1a1a",
      },
    },
  },
  plugins: [],
}
