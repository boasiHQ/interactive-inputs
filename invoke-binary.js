const process = require('node:process');
const os = require('node:os');
const childProcess = require('node:child_process');

const PLATFORM = process.platform;
const CPU_ARCH = os.arch();


function chooseBinary() {
    if ( PLATFORM !== 'linux') {
        throw new Error('Only linux is supported')
    }

    if ( CPU_ARCH !== 'x64' && CPU_ARCH !== 'arm64') {
        throw new Error('Only x64 and arm64 are supported')
    }

    if (PLATFORM === 'linux' && CPU_ARCH === 'x64') {
        return `action-amd64`
    }
    if (PLATFORM === 'linux' && CPU_ARCH === 'arm64') {
        return `action-arm64`
    }
}

const binary = chooseBinary()
const mainScript = `${__dirname}/dist/${binary}`
const spawnSyncReturns = childProcess.spawnSync(mainScript, { stdio: 'inherit' })
process.exit(spawnSyncReturns.status ?? 0)
