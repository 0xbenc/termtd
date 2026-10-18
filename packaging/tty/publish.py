#!/usr/bin/env python3
"""Fetch immutable release bytes or publish them to tty.games. Never rebuild a tag."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import urllib.request

REPO = '0xbenc/termtd'
TARGETS = [(system, arch) for system in ('linux', 'windows', 'darwin') for arch in ('amd64', 'arm64')]


def fetch(url, limit, api=False):
    headers = {'User-Agent': 'termtd-release-publisher'}
    if api and os.environ.get('GH_TOKEN'):
        headers['Authorization'] = 'Bearer ' + os.environ['GH_TOKEN']
    with urllib.request.urlopen(urllib.request.Request(url, headers=headers), timeout=90) as response:
        body = response.read(limit + 1)
    if len(body) > limit:
        raise ValueError('release file exceeds limit')
    return body


def release_files(tag, directory):
    if not re.fullmatch(r'v\d+\.\d+\.\d+(?:-[A-Za-z0-9.-]+)?', tag):
        raise ValueError('explicit release tag required')
    release = json.loads(fetch(f'https://api.github.com/repos/{REPO}/releases/tags/{tag}', 1 << 20, True))
    if release['draft'] or release['tag_name'] != tag:
        raise ValueError('release is not ready')
    assets = {a['name']: a for a in release['assets']}
    directory.mkdir(parents=True, exist_ok=True)

    def download(name, limit):
        asset = assets[name]
        url = f'https://github.com/{REPO}/releases/download/{tag}/{name}'
        if asset['browser_download_url'] != url:
            raise ValueError('unexpected asset origin')
        body = fetch(url, limit)
        sha = hashlib.sha256(body).hexdigest()
        if len(body) != asset['size'] or asset.get('digest') != 'sha256:' + sha:
            raise ValueError('release asset identity mismatch: ' + name)
        return body

    sums = download('checksums.txt', 65536)
    hashes = {}
    for line in sums.decode('ascii').splitlines():
        match = re.fullmatch(r'([a-f0-9]{64})  ([A-Za-z0-9_.-]+)', line)
        if not match or match[2] in hashes:
            raise ValueError('invalid release checksums')
        hashes[match[2]] = match[1]
    for system, arch in TARGETS:
        extension = 'zip' if system == 'windows' else 'tar.gz'
        name = f'termtd_{tag[1:]}_{system}_{arch}.{extension}'
        body = download(name, 128 << 20)
        if hashlib.sha256(body).hexdigest() != hashes.get(name):
            raise ValueError('checksum mismatch: ' + name)
        (directory / name).write_bytes(body)
    (directory / 'checksums.txt').write_bytes(sums)
    print(f'Fetched and verified all six {tag} artifacts; no rebuild')


def publish(tag, directory):
    # This helper/config is checked out at the selected tag for repeatability.
    commit = subprocess.check_output(['git', 'rev-parse', 'HEAD'], text=True).strip()
    tagged = subprocess.check_output(['git', 'rev-parse', tag + '^{commit}'], text=True).strip()
    if commit != tagged:
        raise ValueError('checkout must match the requested release tag')
    pin = json.loads(Path('packaging/tty/tool.json').read_text())
    sha = pin['sha256']
    expected = f'https://tty.games/installers/{sha}/tty-publish-linux-amd64'
    if not re.fullmatch(r'[a-f0-9]{64}', sha) or pin['url'] != expected:
        raise ValueError('invalid publisher pin')
    binary = fetch(expected, 64 << 20)
    if hashlib.sha256(binary).hexdigest() != sha:
        raise ValueError('publisher checksum mismatch')
    tool = directory / 'tty-publish'
    tool.write_bytes(binary)
    tool.chmod(0o700)
    manifest = directory / 'tty-release.json'
    subprocess.run([str(tool.resolve()), 'manifest', '--config', 'tty-release.json', '--version', tag[1:], '--commit', commit, '--dir', str(directory), '--manifest', str(manifest)], check=True)
    endpoint = 'https://tty.games/dev/publish/v1/releases' if '-' in tag else 'https://tty.games/publish/v1/releases'
    subprocess.run([str(tool.resolve()), 'push', '--manifest', str(manifest), '--dir', str(directory), '--endpoint', endpoint], check=True)
    subprocess.run([str(tool.resolve()), 'status', '--game', 'termtd', '--version', tag[1:], '--endpoint', endpoint], check=True)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('operation', choices=['fetch', 'publish'])
    parser.add_argument('--tag', required=True)
    parser.add_argument('--directory', type=Path, default=Path('dist'))
    args = parser.parse_args()
    if args.operation == 'fetch':
        release_files(args.tag, args.directory)
    else:
        publish(args.tag, args.directory)


if __name__ == '__main__':
    main()
