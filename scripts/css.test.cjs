const assert = require('node:assert/strict');
const fs = require('node:fs');
const { test } = require('node:test');
const postcss = require('postcss');

test('compiled CSS supports script-created toast positioning and close transitions', () => {
  const css = postcss.parse(
    fs.readFileSync(
      `${__dirname}/../src/internal/web/ui/static/css/tailwind-base.css`,
      'utf8',
    ),
  );
  function declaration(selector, property) {
    let value;
    css.walkRules((rule) => {
      if (rule.selector.split(',').includes(selector)) {
        rule.walkDecls(property, (decl) => {
          value = decl.value;
        });
      }
    });
    return value;
  }
  assert.equal(declaration('.fixed', 'position'), 'fixed');
  assert.equal(Number(declaration('.opacity-0', 'opacity')), 0);
  assert.equal(Number(declaration('.opacity-30', 'opacity')), 0.3);
  assert.ok(declaration('.max-w-sm', 'max-width'));
});
