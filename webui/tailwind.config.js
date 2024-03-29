const production = !process.env.ROLLUP_WATCH;

module.exports = {
  content: [
    "./src/**/*.svelte",
  ],
  theme: {
    extend: {
      ontFamily: {
        sans: ['Roboto', 'ui-sans-serif', 'system-ui'],
      },
      colors: {
        'primary': 'var(--color-primary)',
        'primary-dark': 'var(--color-primary-dark)',
        'primary-light': 'var(--color-primary-light)'
      },
      textColor: {
        default: 'var(--color-text)',
        standard: 'var(--color-text)'
      }
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
  future: {
    purgeLayersByDefault: true,
    removeDeprecatedGapUtilities: true,
  },
}
