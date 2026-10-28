const assert = require('node:assert/strict');
const { createHash } = require('node:crypto');
const fs = require('node:fs');
const path = require('node:path');

const root = path.resolve(__dirname, '..');
const manifest = JSON.parse(
  fs.readFileSync(path.join(root, 'src/internal/web/ui/vendor.json'), 'utf8'),
);
for (const entry of manifest.packages) {
  for (const asset of [entry, entry.license].filter(Boolean)) {
    const filename = path.resolve(root, asset.path);
    assert.ok(
      filename.startsWith(root + path.sep),
      'Vendor path must stay inside the repository',
    );
    const hash = createHash('sha256')
      .update(fs.readFileSync(filename))
      .digest('hex');
    assert.equal(hash, asset.sha256, `Vendored asset changed: ${asset.path}`);
  }
}
console.log(
  `Verified ${manifest.packages.length} vendored browser packages and their recorded licences`,
);
