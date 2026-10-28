const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const { test } = require('node:test');

const source = fs.readFileSync(`${__dirname}/../invoke-binary.js`, 'utf8');

function invoke(arch, result, platform = 'linux') {
  let exitCode;
  let executable;
  const modules = {
    'node:process': {
      platform,
      exit: (code) => {
        exitCode = code;
      },
    },
    'node:os': { arch: () => arch },
    'node:child_process': {
      spawnSync: (path) => {
        executable = path;
        return result;
      },
    },
  };
  vm.runInNewContext(source, {
    require: (name) => modules[name],
    __dirname: '/action',
    console: { error() {} },
  });
  return { exitCode, executable };
}

test('selects each supported architecture and preserves child exit status', () => {
  assert.deepEqual(invoke('x64', { status: 0 }), {
    exitCode: 0,
    executable: '/action/dist/action-amd64',
  });
  assert.deepEqual(invoke('arm64', { status: 7 }), {
    exitCode: 7,
    executable: '/action/dist/action-arm64',
  });
});

test('failed launches and terminated binaries fail the action', () => {
  assert.equal(
    invoke('arm64', { status: null, error: new Error('ENOENT') }).exitCode,
    1,
  );
  assert.equal(invoke('x64', { status: null, signal: 'SIGTERM' }).exitCode, 1);
});

test('rejects unsupported hosts before launching', () => {
  assert.throws(() => invoke('arm', { status: 0 }), /Only x64 and arm64/);
  assert.throws(() => invoke('arm64', { status: 0 }, 'darwin'), /Only linux/);
});
