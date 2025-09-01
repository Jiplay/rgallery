import * as esbuild from 'esbuild';
import { sassPlugin } from 'esbuild-sass-plugin';

await esbuild.build({
  entryPoints: ['assets/js/main.js', 'assets/js/hls.js', 'assets/styles/main.scss'],
  bundle: true,
  minify: true,
  outdir: 'pkg/dist',
  plugins: [sassPlugin()],
  loader: {
    '.svg': 'dataurl',
    '.png': 'dataurl',
  },
  external: ['*.jpg', '*.woff2'],
  sourcemap: true,
});
