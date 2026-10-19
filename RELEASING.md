# Releasing TERMTD

A stable tag builds one immutable GitHub release, then independently publishes
**tty.games, Homebrew and itch.io**. Current verified stable release: **v1.4.4**.
The workflows, not historical commands below, define the publication behavior.

1. Push a committed `vX.Y.Z` tag.
2. GitHub Actions vets/tests on Linux, macOS and Windows, including packaging tests.
3. GoReleaser builds amd64/arm64 for all three OSes, archives the executables and
   support files, and publishes checksums and SBOMs. Actions attests provenance.
   `GORELEASER_CURRENT_TAG` pins the exact triggering tag.
4. `.github/workflows/publish.yml` fetches the retained release archives, verifies
   both GitHub asset digests and `checksums.txt`, and fans out to three publishers.
   Destination configuration failures do not block the GitHub build or other jobs.
5. itch restages Play bundles from those exact binaries; Homebrew updates
   `Formula/termtd.rb` in `0xbenc/homebrew-tap`; tty.games verifies and stores the
   immutable archives before activating its stable channel.

Install through Homebrew with `brew install 0xbenc/tap/termtd`.
Prerelease tags publish GitHub assets and tty.games **testing** only; stable
Homebrew and itch channels remain unchanged.

## tty.games

`tty-release.json` declares the command, support files and reviewed save contract.
`packaging/tty/tool.json` pins the shared publisher's URL and SHA-256. The helper
uses `/publish/v1/releases` for stable and `/dev/publish/v1/releases` for testing.
GitHub OIDC authenticates the repository, workflow, source commit and tag; no AWS
key or long-lived tty.games secret is stored in this game's Actions settings.

The deployed server policy must authorize this game and its exact install/save
contract and platform set. Server-side AWS permissions provide scoped immutable
S3 uploads. Identical retries succeed; conflicting releases and backward channel
updates fail. Players explicitly install/update through the storefront; updates
preserve the reviewed save contract and back up saves.

## Retry an existing release

In Actions, run **Publish existing release** from `main`, supply the existing tag,
and choose `tty`, `homebrew`, `itch`, or `all`. Recovery downloads the retained
release; it does not rebuild or require a new tag. Homebrew rejects downgrades;
itch requires the latest stable GitHub release. Support begins at v1.4.3-rc.1.

Verified example: v1.4.4 release run 39098534082 succeeded for the build, tty.games
and Homebrew. Its itch job failed after the account rename. The itch-only retry
39098967683 succeeded, with the other destinations skipped. Current itch target:
`tty-games/termtd`. Do not interpret the original run's failure as all destinations
having failed.

## itch packages

The itch destination is <https://tty-games.itch.io/termtd>. Each build has an
explicit `.itch.toml` Play action and instructions in `README-PLAY.txt`:

| Channel | Player's launcher |
| --- | --- |
| `linux-amd64` | `Play.sh` → desktop terminal, or current terminal |
| `linux-arm64` | `Play.sh` → desktop terminal, or current terminal |
| `mac-amd64` | `TERMTD.app` → Terminal; `Play.command` also works |
| `mac-arm64` | `TERMTD.app` → Terminal; `Play.command` also works |
| `windows-amd64` | `termtd.exe` with itch's `console = true`; `Play.cmd` for browser downloads |
| `windows-arm64` | `termtd.exe` with itch's `console = true`; `Play.cmd` for browser downloads |

Launchers do not install anything or change PATH. Homebrew provides the global
`termtd` command. The macOS bundle requires macOS 12+, matching
[Go 1.26's platform support](https://go.dev/doc/go1.26#darwin).

Butler is pinned to 15.31.0 with an archive SHA-256 in
`packaging/itch/install-butler.sh`. The builder uses GoReleaser's
`dist/artifacts.json` to locate the binaries and restores executable permissions
when rebuilding packages from downloaded CI artifacts.

Windows packages require Windows 10/11. The executable uses native console
keyboard, mouse, and resize events and enables UTF-8 virtual terminal output.
The itch manifest opens a console directly; Windows Terminal is optional.
GitHub Windows archives and portable itch downloads use ZIP format.

Itch publishing validates all six packages before uploading the first one,
then pushes with `--userversion X.Y.Z`. A retry replaces the same channels;
publishing across all six channels is sequential, not an atomic transaction.

## One-time GitHub setup

The TERMTD repository needs an Actions secret named `TAP_GITHUB_TOKEN` with
Contents write access to `0xbenc/homebrew-tap`. This is the same secret name
used by the other projects. TERMTD's secret has been configured using the existing
local `gh` login, which has tap write access and TERMTD admin access.

For another repository, check the existing login before requesting a token:

```sh
gh auth status
gh api repos/0xbenc/homebrew-tap --jq '.permissions.push'
gh api repos/0xbenc/termtd --jq '.permissions.admin'
```

When that login has the required access, configure the secret directly without
printing or saving its value:

```sh
set -o pipefail
gh auth token --hostname github.com | gh secret set TAP_GITHUB_TOKEN --repo 0xbenc/termtd
```

GitHub does not return existing repository secret values. Reusing the local
authenticated credential avoids needing to retrieve another repository's secret.
Alternatively, supply a dedicated tap token through the private prompt:

```sh
gh secret set TAP_GITHUB_TOKEN --repo 0xbenc/termtd
```

A missing tap secret fails only the Homebrew publishing job; the GitHub release
and other destinations remain independent.
The regular `GITHUB_TOKEN` is supplied
by Actions for publishing the TERMTD release itself.

The following itch settings are configured on `0xbenc/termtd`:

- Actions variable `ITCH_TARGET`: `tty-games/termtd`.
- Actions secret `BUTLER_API_KEY`: copied directly from the local Butler login.

If the itch credential needs renewal, run `butler login` locally, then transfer
its saved value without printing it:

```sh
gh secret set BUTLER_API_KEY --repo 0xbenc/termtd < ~/.config/itch/butler_creds
```

See [Butler authentication](https://itch.io/docs/butler/login.html) for the
credential location on other OSes. Each destination checks its own credentials when its publishing job runs.
Missing itch credentials do not block the GitHub release, tty.games or Homebrew.

The setup commands above document credential provisioning; they do not prove
credentials are available on every workstation. Never print credential values.

## Prepare a release

The currently verified stable release is 1.4.4. Use a new, unused version consistently in the
release tag, itch user version, release notes, and post-release checks. Examples
below use `VERSION=1.4.5` as the next-version example; check tags before choosing it. Prerelease tags
publish GitHub assets and tty.games testing without updating stable itch or Homebrew channels.

Player archives contain LICENSE, README.md, LORE.md, and PLAYING.md. The lore
summary follows the current opening and avoids ending spoilers. Development
balance reports remain in TACTICS.md in the repository.

Review and commit the intended game changes and release configuration first.
Generated images and release artifacts under `output/` and `dist/` are ignored.

```sh
go vet ./...
go test ./...
python3 -m unittest discover -s packaging/homebrew -p 'test_*.py'
python3 -m unittest discover -s packaging/itch -p 'test_*.py'
goreleaser check
goreleaser release --snapshot --clean --skip=sbom
python3 packaging/itch/stage.py
butler validate --platform linux --arch amd64 dist/itch/linux-amd64
butler validate --platform osx --arch amd64 dist/itch/mac-amd64
```

The snapshot command builds all six targets and archives without publishing.
Remove `--skip=sbom` if `syft` is installed to also exercise SBOM generation.
Release CI installs syft before publishing. Local binaries report `termtd dev`;
GoReleaser embeds the version from the tag.

Before tagging, confirm that Linux, macOS, and Windows CI passed and smoke-test
playing, resizing, mouse input, exiting, and continuing a saved campaign on all
three OSes.
Cross-compilation checks builds, while a real terminal checks platform behavior.
Also test the macOS app and `.command` launcher and the Linux graphical Play
launcher from the itch app. Linux launcher tests cover terminal selection and
paths/arguments with spaces; a graphical macOS launch requires a Mac.
Butler's validator currently accepts only `386` and `amd64`, so ARM package
validation checks the script/app launch target using `--arch amd64`.

Windows ARM packages contain a native ARM executable. The pinned Butler validates
their launch manifest but warns about the ARM PE header and reports amd64, so its
output does not verify their executable architecture. Smoke-test on ARM Windows
before release.
On Windows, also test itch Play and Play.cmd from an extracted folder with spaces,
both in Windows Terminal and the classic console. Verify that normal exit restores
the console settings and startup errors remain visible through Play.cmd.

To preview the stable channel commands using staged stable-version packages:

```sh
VERSION=1.4.5
python3 packaging/itch/publish.py --target tty-games/termtd --version "$VERSION" --dry-run
```

Remove `--dry-run` only when intentionally publishing those staged packages.

When ready to publish the committed revision:

```sh
VERSION=1.4.5
git tag -a "v$VERSION" -m "TERMTD $VERSION"
git push origin main
git push origin "v$VERSION"
```

After the release workflow succeeds:

```sh
brew update
brew install 0xbenc/tap/termtd
termtd --version
brew test 0xbenc/tap/termtd
butler status tty-games/termtd
```

If the GitHub release succeeds but the tap push fails, use its existing
`checksums.txt` with `packaging/homebrew/render.py` to regenerate the formula;
do not rebuild the published archives with different checksums. The publisher
leaves the tap unchanged if that formula is already current.

## Save compatibility

Player data uses `~/.termtd-hiscores.json`, `~/.termtd-journal.json`, and
`~/.termtd-lair.json`. If a new file is absent, the game reads its old `.tdef-`
counterpart. Subsequent saves use the new name and leave old files intact.
Reset progress removes both sets so legacy data cannot reappear.

These files record campaign results, roster unlocks, discoveries, relic tokens,
and scores. They do not resume a defense in progress. Save errors are shown in
the game, retained in memory, and retried automatically or with Ctrl+L.

Existing saves that cannot be read or decoded now produce a persistent load
warning. Missing files are normal first runs; invalid JSON (including a top-level
null or an incompatible field type) is a load failure, with partial decoded data
discarded. Corrupt files are copied byte-for-byte to a unique sibling
`.corrupt-*` file, with copy failures reported, and originals stay untouched.
A failed current save never falls back to a legacy save.

The affected save is blocked for the session, including automatic writes, retries,
and screen refreshes. Ctrl+L opens scrollable recovery details and pauses a running
defense. Repair/restore the reported file and restart, or explicitly confirm Reset
progress. Reset clears the load block but retains recovery copies. Other healthy
saves can still be written; temporary progress is not merged into restored data.

Before releasing, test a malformed campaign/journal/score file and an unreadable
save path. Verify the warning and recovery bytes, quit confirmation, repair plus
restart, and confirmed reset. Repeat with only legacy saves present.
