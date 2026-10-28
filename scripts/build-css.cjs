const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const { execFileSync } = require('node:child_process');

const root = path.resolve(__dirname, '..');
const args = process.argv.slice(2);
if (args.some((arg) => arg !== '--check')) {
  throw new Error('Usage: node scripts/build-css.cjs [--check]');
}
const checking = args.includes('--check');
const target = path.join(
  root,
  'src/internal/web/ui/static/css/tailwind-base.css',
);
const temporary = checking
  ? fs.mkdtempSync(path.join(os.tmpdir(), 'interactive-inputs-css-'))
  : null;
const output = checking ? path.join(temporary, 'tailwind.css') : target;
try {
  execFileSync(
    process.execPath,
    [
      require.resolve('tailwindcss/lib/cli.js'),
      '--config',
      path.join(root, 'tailwind.config.cjs'),
      '--input',
      path.join(root, 'src/styles/tailwind.css'),
      '--output',
      output,
      '--minify',
    ],
    { cwd: root, stdio: 'inherit' },
  );
  if (checking && !fs.readFileSync(target).equals(fs.readFileSync(output))) {
    throw new Error(
      'Generated CSS is stale; run yarn build:css and yarn package',
    );
  }
} finally {
  if (temporary) fs.rmSync(temporary, { recursive: true, force: true });
}
