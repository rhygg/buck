const os = require('os')
const json = require('../package.json')
const { downloadRelease } = require('@terascope/fetch-github-release')
const path = require('path')
const { extract } = require('tar')
const { copy, remove, writeJson } = require('fs-extra')

const archMap = {
  arm64: 'arm64',
  x64: 'amd64',
}
const platformMap = {
  win32: 'windows',
  linux: 'linux',
  darwin: 'macos',
}

async function main() {
  const platform = platformMap[os.platform()]
  const arch = os.arch()
  const tempPath = path.resolve(__dirname, '.temp')
  const ext = platform === 'windows' ? '.exe' : ''
  const binPath = path.resolve(__dirname, 'bin' + ext)
  const version = json.version
  const assetName = `buck_${version}_${platform}_${archMap[arch]}.tar.gz`
  try {
    await downloadRelease(
      'rhygg',
      'buck',
      tempPath,
      (release) => {
        return release.tag_name === `v${version}`
      },
      (asset) => {
        return asset.name === assetName
      },
      false,
      false,
    )
  } catch (e) {
    console.error('There seems to be an error while downloading: ', e)
  }
  await extract({
    file: path.resolve(tempPath, assetName),
    cwd: tempPath,
  })

  await copy(path.resolve(tempPath, 'buck' + ext), binPath)
  await remove(path.resolve(__dirname, 'bin'))
  json.bin.saki = 'bin' + ext
  await writeJson(path.resolve(__dirname, 'package.json'), json)
}

main()