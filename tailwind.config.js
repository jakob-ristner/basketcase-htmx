/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./internal/**/*.{go,js,templ,html}"
  ],
  theme: {
    fontFamily: {
      'mono': ['Consolas'],
      'display': ['Playfair Display'],
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
  plugins: [
    require("@catppuccin/tailwindcss")({
      // prefix to use, e.g. `text-pink` becomes `text-ctp-pink`.
      // default is `false`, which means no prefix
      prefix: "ctp",
      // which flavour of colours to use by default, in the `:root`
      defaultFlavour: "latte",
    }),
  ],
}