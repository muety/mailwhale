import svelte from 'rollup-plugin-svelte'
import sveltePreprocess from 'svelte-preprocess'
import commonjs from '@rollup/plugin-commonjs'
import resolve from '@rollup/plugin-node-resolve'
import livereload from 'rollup-plugin-livereload'
import { terser } from 'rollup-plugin-terser'
import css from 'rollup-plugin-css-only'
import replace from '@rollup/plugin-replace'
import dev from 'rollup-plugin-dev'

const production = !process.env.ROLLUP_WATCH

export default {
  input: 'src/main.js',
  output: {
    sourcemap: true,
    format: 'iife',
    name: 'app',
    file: 'public/build/bundle.js'
  },
  plugins: [
    svelte({
      compilerOptions: {
        // enable run-time checks when not in production
        dev: !production
      },
      preprocess: sveltePreprocess({
        sourceMap: !production,
        postcss: {
          plugins: [
            require('tailwindcss'),
            require('autoprefixer'),
          ],
        },
      }),
      onwarn: (warning, defaultHandler) => {
        if (warning.code === 'a11y-distracting-elements') return;
        if (warning.code === 'a11y-click-events-have-key-events') return;
        if (warning.code === 'css-unused-selector') return;

        // handle all other warnings normally
        defaultHandler(warning);
      },
    }),
    // we'll extract any component CSS out into
    // a separate file - better for performance
    css({ output: 'bundle.css' }),

    replace({
      preventAssignment: true,
      process: JSON.stringify({
        env: {
          isDev: !production,
          apiUrl: process.env.API_URL,
        }
      }),
    }),

    // If you have external dependencies installed from
    // npm, you'll most likely need these plugins. In
    // some cases you'll need additional configuration -
    // consult the documentation for details:
    // https://github.com/rollup/plugins/tree/master/packages/commonjs
    resolve({
      browser: true,
      dedupe: ['svelte']
    }),
    commonjs(),

    // In dev mode, call `npm run start` once
    // the bundle has been generated
    !production && dev({
      dirs: ['public'],
      proxy: [{ from: '/api/*', to: 'http://127.0.0.1:3000/' }],
      port: 5000,
      spa: true
    }),

    // Watch the `public` directory and refresh the
    // browser on changes when not in production
    !production && livereload('public'),

    // If we're building for production (npm run build
    // instead of npm run dev), minify
    production && terser()
  ],
  watch: {
    clearScreen: false
  }
}
