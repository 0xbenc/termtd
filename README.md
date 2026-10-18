# termtd

Tower defense in your terminal.

You’re Grak, the Last Monster. Heroes are coming for the hoard, and the dragon
who lives on it is dying. Build a defense. Keep him alive.

Eleven defenders, twenty waves, a lair to explore, and a journal full of lore.
Written in Go. Zero dependencies.

## Get it

Download it from [itch.io](https://kairuku-studios.itch.io/termtd) and use the Play
launcher, or install it from my Homebrew tap:

```sh
brew install 0xbenc/tap/termtd
termtd
```

For Windows 10/11, Linux, and macOS (12+), on Intel/AMD and ARM.

## Build it

With Go 1.26.3 or newer:

```sh
go build -o termtd .
./termtd
```

On Windows, build and run from PowerShell:

```powershell
go build -o termtd.exe .
.\termtd.exe
```

## Play

Use the arrows or WASD to move. Pick a tower with `1`–`7`, then press Enter
or click to place it. Green dots mark sites you can afford; amber dots mark sites
that need more gold. Each build exits placement; pick a tower again to build
another. Escape cancels placement.

- `[` / `]` — switch between the original defenders and four specialists
- `r` — rotate a Runeforge while placing it or after selecting it
- `u` / `x` — upgrade / sell
- `Tab` — cycle through placed defenders; moving onto one also shows its controls
- `t` — cycle targeting priority; the dotted guide shows its current target
- `n` — send the next wave
- `p` / `f` — pause / change speed
- `j` — open the journal in the lair or while paused
- `h` — help · `esc` — back · `q` — quit

Each specialist is limited to one placed per level. Selling one frees its slot.

A fresh campaign teaches the main roster during the first Rotunda defense,
then introduces the specialists during the second defense. Recruitment pauses
the action and supplies gold for the new defender. Placing a defender resumes
the automatic wave countdown; `p` pauses and `n` sends a wave early. Earned
recruits stay unlocked on retries. Quick Play keeps the full roster.

Replay either guided defense without changing your saves:

```powershell
.\termtd.exe tutorial
.\termtd.exe tutorial -stage 2
```

Completed defenses, recruits, discoveries, and scores save automatically.
An unfinished defense restarts when you return. Use a terminal of at least
62 columns by 19 rows; more room gives the artwork and controls extra space.

**Start** advances your campaign in the lair. **Quick Play** lets you choose any
regular map or a procedural maze and records high scores, without changing
campaign progress, relics, journal discoveries, or story progress.

Open **Journal** from the main menu to browse discoveries earned in the campaign.

To start over, choose **Reset progress** in the main menu, then confirm
**Reset all progress**. This clears your campaign, relics, journal, and high scores.

[Story](LORE.md) · [Lore sources](https://github.com/0xbenc/termtd/blob/main/LORE-SOURCES.md) · [Field guide](PLAYING.md) · [Release setup](https://github.com/0xbenc/termtd/blob/main/RELEASING.md) · [MIT license](LICENSE)

### Inspecting saves and updating

Run `termtd saves` to see the save locations and check whether this build can
read them, or `termtd saves --json` for a machine-readable report. This command
never rewrites, repairs or migrates files. Missing files are normal on a first
run; unreadable or invalid files produce a nonzero exit status. The existing
save format remains `termtd-v1` in 1.4.3.

Update with the tool that installed your copy: Homebrew, itch, or the terminal
installer at `ssh tty.games`. Close the game before updating. tty.games asks
before replacing its own executable and backs up existing saves; it does not
update a Homebrew or itch installation. Keep backups before changing versions.

Release publishing builds once and sends the retained artifacts independently
to Homebrew, itch and tty.games. Stable tags publish after tests and verification;
prereleases go only to tty.games testing. To recover a destination failure, use
**Actions → Publish existing release → Run workflow**, select `main`, enter the
existing tag and choose the destination. No rebuild or new tag is needed. The
first tag supporting this recovery workflow is v1.4.3-rc.1. Publishing an older
version cannot move tty.games stable backward.
