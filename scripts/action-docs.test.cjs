const assert = require('node:assert/strict');
const { test } = require('node:test');
const { generate } = require('./action-docs.cjs');

const marker = '<!-- action-docs-inputs source="action.yml" -->';
const metadata = {
  name: 'Example',
  description: 'Example action',
  runs: { using: 'node24' },
  inputs: {
    fields: {
      description: 'A|B\n<script>`',
      required: true,
      default: 'fields:\n  - label: test',
    },
    optional: { default: false },
  },
};

test('regenerates only input markers, escapes metadata and retains prose', () => {
  const before = `Intro\n${marker}\nstale\n${marker}\nUsage`;
  const result = generate(before, metadata);
  assert.ok(result.startsWith('Intro\n'));
  assert.ok(result.endsWith('\nUsage'));
  assert.match(result, /A&#124;B<br>&lt;script&gt;&#96;/);
  assert.match(result, /`fields:<br>  - label: test`/);
  assert.match(result, /\| `optional` \| <p><\/p> \| `false` \| `false` \|/);
  assert.equal(generate(result, metadata), result);
});

test('fails on incomplete markers and quoted required booleans', () => {
  assert.throws(() => generate('', metadata), /exactly two inputs markers/);
  assert.throws(
    () =>
      generate(`${marker}\n${marker}`, {
        ...metadata,
        inputs: { field: { required: 'false' } },
      }),
    /required must be a YAML boolean/,
  );
});
