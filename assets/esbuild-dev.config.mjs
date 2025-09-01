import * as esbuild from 'esbuild';
import { sassPlugin } from 'esbuild-sass-plugin';
import { livereloadPlugin } from '@jgoz/esbuild-plugin-livereload';

let ctx = await esbuild.context({
  entryPoints: ['assets/js/main.js', 'assets/js/hls.js', 'assets/styles/main.scss'],
  outdir: 'pkg/dist',
  bundle: true,
  plugins: [sassPlugin(), livereloadPlugin()],
  loader: {
    '.svg': 'dataurl',
    '.png': 'dataurl',
  },
  external: ['*.jpg', '*.woff2'],
  sourcemap: true,
});

await ctx.watch();

let { host, port } = await ctx.serve({
  servedir: 'pkg/dist',
});
