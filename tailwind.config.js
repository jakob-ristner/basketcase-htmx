/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./internal/**/*.{go,js,templ,html}"
  ],
  theme: {
    extend: {
      colors: {
        // 'regal-blue': '#243c5a',
        // 'light': '#F2F1EB',
        // 'cream': '#EEE7DA',
        // 'l-green': '#AFC8AD',
        'm-green': '#88AB8E',

        'dark-moss-green': '#606c38',
        'pakistan-green': '#283618',
        'cornsilk': '#fefae0',
        'earth-yellow': '#dda15e',
        'tigers-eye': '#bc6c25',
        'dark-text': '#2D2424',
      },
    },
    fontFamily: {
      //'sans': ['ui-sans-serif', 'system-ui'],
      //'serif': ['Playfair Display'],
      //'mono': ['ui-monospace', 'SFMono-Regular'],
      'display': ['Playfair Display'],
      //'body': ['"Open Sans"'],
    },
    fontSize: {
      sm: '0.8rem',
      base: '1rem',
      xl: '1.25rem',
      '2xl': '2rem',
      '3xl': '3rem',
      '4xl': '4rem',
      '5xl': '5rem',
    },
  },
  plugins: [],
}