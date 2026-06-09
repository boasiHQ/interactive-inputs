const process = require('node:process');
const os = require('node:os');
const childProcess = require('node:child_process');
const path = require('node:path');

const PLATFORM = process.platform;
const CPU_ARCH = os.arch();

/**
 * Selects the prebuilt action binary for the current GitHub Actions runner.
 *
 * Interactive Inputs ships Linux binaries because JavaScript actions execute
 * on the runner where the workflow job is already running.
 *
 * @returns {'action-amd64' | 'action-arm64'} The binary name under dist/.
 * @throws {Error} When the runner platform or CPU architecture is unsupported.
 */
function chooseBinary() {
  if (PLATFORM !== 'linux') {
    throw new Error('Only linux is supported');
  }

  if (CPU_ARCH !== 'x64' && CPU_ARCH !== 'arm64') {
    throw new Error('Only x64 and arm64 are supported');
  }

  if (CPU_ARCH === 'x64') {
    return 'action-amd64';
  }

  return 'action-arm64';
}

/**
 * Runs the selected Go action binary and mirrors its exit state back to the
 * JavaScript action runtime.
 *
 * @param {string} binary The binary filename returned by {@link chooseBinary}.
 */
function invokeBinary(binary) {
  const mainScript = path.join(__dirname, 'dist', binary);
  const spawnSyncReturns = childProcess.spawnSync(mainScript, {
    stdio: 'inherit',
  });

  if (spawnSyncReturns.error) {
    throw spawnSyncReturns.error;
  }

  if (spawnSyncReturns.signal) {
    console.error(
      `Action binary exited due to signal ${spawnSyncReturns.signal}`,
    );
    process.exit(1);
  }

  process.exit(spawnSyncReturns.status ?? 1);
}

invokeBinary(chooseBinary());
