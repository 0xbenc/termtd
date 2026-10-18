package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/0xbenc/termtd/game"
	"github.com/0xbenc/termtd/hiscore"
	"github.com/0xbenc/termtd/internal/buildinfo"
	"github.com/0xbenc/termtd/render"
	"github.com/0xbenc/termtd/tui"
)

// Release builds set version through -ldflags; local builds remain identifiable.
var version = "dev"

func main() {
	buildinfo.Version = version
	if len(os.Args) < 2 {
		play(nil)
		return
	}
	switch os.Args[1] {
	case "version", "--version", "-version":
		fmt.Printf("termtd %s\n", version)
	case "saves":
		flags := flag.NewFlagSet("saves", flag.ExitOnError)
		asJSON := flags.Bool("json", false, "print a machine-readable save report")
		flags.Parse(os.Args[2:])
		report, err := hiscore.Inspect()
		if err != nil {
			die("%v", err)
		}
		if err = hiscore.WriteReport(os.Stdout, report, *asJSON); err != nil {
			die("%v", err)
		}
		if !report.Healthy {
			os.Exit(1)
		}
	case "play":
		play(os.Args[2:])
	case "bench":
		bench(os.Args[2:])
	case "headless":
		headless(os.Args[2:])
	case "capture":
		capture(os.Args[2:])
	case "overworld":
		overworldCmd()
	case "intro", "ending":
		film := render.FilmOpening
		if os.Args[1] == "ending" {
			film = render.FilmEnding
		}
		if err := tui.RunCutscene(film); err != nil {
			die("%v", err)
		}
	case "tutorial":
		flags := flag.NewFlagSet("tutorial", flag.ExitOnError)
		stage := flags.Int("stage", 1, "guided defense: 1 Rotunda, 2 specialists in the Rift")
		flags.Parse(os.Args[2:])
		if err := tui.RunTutorial(*stage); err != nil {
			die("%v", err)
		}
	case "journal":
		if err := tui.RunJournal(); err != nil {
			die("%v", err)
		}
	case "centurion":
		if err := tui.RunCenturionMockup(); err != nil {
			die("%v", err)
		}
	case "squire":
		if err := tui.RunSquireMockup(); err != nil {
			die("%v", err)
		}
	case "wizard":
		if err := tui.RunWizardMockup(); err != nil {
			die("%v", err)
		}
	case "mercenary":
		if err := tui.RunMercenaryMockup(); err != nil {
			die("%v", err)
		}
	case "rogue":
		if err := tui.RunRogueMockup(); err != nil {
			die("%v", err)
		}
	case "paladin":
		if err := tui.RunPaladinMockup(); err != nil {
			die("%v", err)
		}
	case "necromancer":
		if err := tui.RunNecromancerMockup(); err != nil {
			die("%v", err)
		}
	case "trebuchet":
		if err := tui.RunTrebuchetMockup(); err != nil {
			die("%v", err)
		}
	case "lightning":
		if err := tui.RunLightningMockup(); err != nil {
			die("%v", err)
		}
	case "ranger":
		if err := tui.RunRangerMockup(); err != nil {
			die("%v", err)
		}
	case "cannonier":
		if err := tui.RunCannonierMockup(); err != nil {
			die("%v", err)
		}
	case "frost":
		if err := tui.RunFrostMockup(); err != nil {
			die("%v", err)
		}
	case "player":
		if err := tui.RunPlayerMockup(); err != nil {
			die("%v", err)
		}
	case "gunner":
		if err := tui.RunGunnerMockup(); err != nil {
			die("%v", err)
		}
	case "gnolls":
		if err := tui.RunGnollMockup(); err != nil {
			die("%v", err)
		}
	case "grak":
		if err := tui.RunGrakMockup(); err != nil {
			die("%v", err)
		}
	case "dragon":
		if err := tui.RunDragonMockup(); err != nil {
			die("%v", err)
		}
	case "maps":
		mapsCmd()
	case "help", "-h", "--help":
		usage()
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `termtd - terminal tower defense

usage:
  termtd [play flags]        play (default: title -> menu -> the lair)
  termtd saves [--json]      inspect save locations and readability without changes
  termtd bench [flags]       run autoplay balance benchmark
  termtd headless [flags]    run one autoplay game, print result
  termtd capture [flags]     render headless game frames to files
  termtd overworld           the lair map: walk the floors, descend into one
  termtd intro               replay the opening film
  termtd ending              replay the Heart ending film
  termtd journal             browse Grak's lore journal
  termtd tutorial [-stage 2]  replay a guided defense without changing saves
  termtd dragon              view Malgrath's static portrait mockup
  termtd grak                view Grak's static portrait mockup
  termtd gnolls              view the Gnoll Slingers lore portrait
  termtd gunner              view the Orc Gunner lore portrait
  termtd frost               view the Frost Mage lore portrait
  termtd player              view the guild champion lore portrait
  termtd centurion           view the warded guild veteran portrait
  termtd squire              view the young guild recruit portrait
  termtd wizard              view the guild spellcaster portrait
  termtd mercenary           view the weary guild hireling portrait
  termtd rogue               view the guild knife-runner portrait
  termtd paladin             view the guild vanguard portrait
  termtd necromancer         view the guild summoner portrait
  termtd trebuchet           view the monster siege crew portrait
  termtd lightning           view the Lightning Mage portrait
  termtd ranger              view the monster archer portrait
  termtd cannonier           view the monster artillery crew portrait
  termtd maps                list built-in levels
  termtd version             print the binary version
  termtd help

play flags:
  -level string   level name (default "winding"; see termtd maps)
  -maze int       use a procedural maze with this seed (0 = random)

bench flags:
  -n int          games per level (default 40)
  -levels string  comma list of levels, or "all" (default "all")
  -maze int       include N procedural mazes (default 40)

capture flags:
  -out string     output dir (default "frames")
  -every int      render every N ticks (default 10)
  -limit int      stop after N ticks (default: game end)
  -text           write plain text instead of ANSI
`)
}

func parseLevelArgs(fs *flag.FlagSet, args []string) (*game.Map, string, *flag.FlagSet) {
	level := fs.String("level", "winding", "")
	mazeSeed := fs.Int64("maze", -1, "")
	_ = fs.Parse(args)
	if *mazeSeed >= 0 {
		// -maze 0 means "random maze": resolve the seed here (and name the
		// level after the real seed, so hiscores stay per-maze).
		ms := pickSeed(*mazeSeed)
		m, err := game.MazeFromSeed(ms)
		if err != nil {
			die("maze: %v", err)
		}
		return m, fmt.Sprintf("maze%d", ms), fs
	}
	m, err := game.LoadLevel(*level)
	if err != nil {
		die("%v", err)
	}
	return m, *level, fs
}

func pickSeed(seed int64) int64 {
	if seed != 0 {
		return seed
	}
	return time.Now().UnixNano()
}

func parseDiff(fs *flag.FlagSet) *flag.FlagSet {
	fs.String("diff", "normal", "difficulty: easy|normal|hard")
	return fs
}

func diffFrom(fs *flag.FlagSet) game.Difficulty {
	s := fs.Lookup("diff").Value.String()
	switch s {
	case "easy":
		return game.Easy
	case "hard":
		return game.Hard
	}
	return game.Normal
}

func play(args []string) {
	fs := flag.NewFlagSet("play", flag.ExitOnError)
	fs.String("diff", "easy", "difficulty: easy|normal|hard")
	m, name, _ := parseLevelArgs(fs, args)
	diff := diffFrom(fs)
	explicit := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "level" || f.Name == "maze" {
			explicit = true
		}
	})
	if explicit {
		if err := tui.Run(m, name, diff); err != nil {
			die("%v", err)
		}
		return
	}
	// No level on the command line: start at the title screen and let
	// the player pick a level, maze seed and difficulty.
	if err := tui.RunMenu(diff); err != nil {
		die("%v", err)
	}
}

func bench(args []string) {
	fs := flag.NewFlagSet("bench", flag.ExitOnError)
	parseDiff(fs)
	n := fs.Int("n", 40, "")
	levels := fs.String("levels", "all", "")
	nMaze := fs.Int("maze", 40, "")
	_ = fs.Parse(args)
	diff := diffFrom(fs)

	names := game.LevelNames()
	if *levels == "all" {
		// The heart is the endgame gate; bench it like a level.
		names = append(names, "heart")
	} else {
		names = strings.Split(*levels, ",")
	}
	total := len(names)*(*n) + *nMaze
	fmt.Printf("bench: %d games\n", total)
	start := time.Now()
	all := []game.SimResult{}
	for _, name := range names {
		var m *game.Map
		var err error
		if name == "heart" {
			m, err = game.LoadBoss()
		} else {
			m, err = game.LoadLevel(name)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "skip:", err)
			continue
		}
		for i := 0; i < *n; i++ {
			all = append(all, game.RunAutoplayDiff(m, name, diff))
		}
	}
	for i := 0; i < *nMaze; i++ {
		s := int64(100000 + i)
		m, err := game.MazeFromSeed(s)
		if err != nil {
			continue
		}
		all = append(all, game.RunAutoplayDiff(m, "maze", diff))
	}
	byLevel := map[string][]game.SimResult{}
	for _, r := range all {
		byLevel[r.Level] = append(byLevel[r.Level], r)
	}
	fmt.Printf("%-10s %8s %10s %10s %8s %8s %8s\n", "level", "winrate", "avgWave", "avgLeaks", "towers", "gold", "time")
	wins, count := 0, 0
	for _, name := range append(append([]string{}, names...), "maze") {
		rs, ok := byLevel[name]
		if !ok || len(rs) == 0 {
			continue
		}
		w, waves, leaks, towers, gold, t := stats(rs)
		fmt.Printf("%-10s %7.0f%% %10.1f %10.2f %8.1f %8.0f %8.0f\n",
			name, 100*float64(w)/float64(len(rs)), waves, leaks, towers, gold, t)
		wins += w
		count += len(rs)
	}
	if count > 0 {
		fmt.Printf("%-10s %7.0f%%\n", "ALL", 100*float64(wins)/float64(count))
	} else {
		fmt.Printf("%-10s %7s\n", "ALL", "n/a")
	}
	fmt.Printf("elapsed %v\n", time.Since(start).Round(time.Millisecond))
}

func stats(rs []game.SimResult) (wins int, waves, leaks, towers, gold, t float64) {
	for _, r := range rs {
		if r.Won {
			wins++
		}
		waves += float64(r.Wave)
		// r.Leaks counts actual leaks. 20-r.Lives would overcount boss
		// leaks (6 lives each) and be wrong on hard (15 starting lives).
		leaks += float64(r.Leaks)
		towers += float64(r.Towers)
		gold += float64(r.Gold)
		t += r.Time
	}
	n := float64(len(rs))
	waves, leaks, towers, gold, t = waves/n, leaks/n, towers/n, gold/n, t/n
	return
}

func headless(args []string) {
	fs := flag.NewFlagSet("headless", flag.ExitOnError)
	parseDiff(fs)
	m, name, _ := parseLevelArgs(fs, args)
	r := game.RunAutoplayDiff(m, name, diffFrom(fs))
	fmt.Printf("%s: %s wave=%d lives=%d gold=%d towers=%d kills=%d time=%.0fs\n",
		name, status(r.Won), r.Wave, r.Lives, r.Gold, r.Towers, r.Kills, r.Time)
}

func status(won bool) string {
	if won {
		return "WIN"
	}
	return "LOSE"
}

func capture(args []string) {
	fs := flag.NewFlagSet("capture", flag.ExitOnError)
	parseDiff(fs)
	out := fs.String("out", "frames", "")
	every := fs.Int("every", 10, "")
	limit := fs.Int64("limit", 0, "")
	text := fs.Bool("text", false, "")
	scale := fs.Int("scale", 1, "playfield scale 1-4")
	m, name, _ := parseLevelArgs(fs, args)
	if *every < 1 {
		die("capture: -every must be >= 1")
	}
	if *scale < 1 || *scale > 4 {
		die("capture: -scale must be 1-4")
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		die("%v", err)
	}
	// The virtual terminal is pinned to the minimum size that yields the
	// requested scale, so output is identical in any environment (the
	// real tty size is never inherited).
	tw, th := render.CaptureSize(m.W, m.H, *scale)
	s := game.NewStateDiff(m, diffFrom(fs))
	ai := game.NewAutoplay(s)
	pal := render.Palette()
	ui := render.UI{Cursor: game.Vec{X: m.W / 2, Y: m.H / 2}, Placing: game.TowerGunner, Selected: render.NoSelection, Speed: 1, Level: name}
	dt := 1.0 / 20.0
	tick := int64(0)
	for s.Status == game.StatusRunning {
		ai.Tick()
		s.Step(dt)
		tick++
		if tick%int64(*every) == 0 {
			f := render.Render(s, &ui, pal, tw, th, int(tick))
			p := filepath.Join(*out, fmt.Sprintf("%s_%06d", name, tick))
			var content string
			if *text {
				content = f.Text() + "\n"
			} else {
				content = f.ANSI() + "\n"
			}
			os.WriteFile(p, []byte(content), 0o644)
		}
		if *limit > 0 && tick >= *limit {
			break
		}
	}
	st := "running"
	if s.Status == game.StatusVictory {
		st = "WIN"
	} else if s.Status == game.StatusDefeat {
		st = "LOSE"
	}
	fmt.Printf("wrote frames to %s (ticks %d, wave %d, %s)\n", *out, tick, s.Wave, st)
}

func overworldCmd() {
	if err := tui.RunOverworld(); err != nil {
		die("%v", err)
	}
}

func mapsCmd() {
	fmt.Println("built-in levels:")
	for _, n := range game.LevelNames() {
		m, err := game.LoadLevel(n)
		if err != nil {
			fmt.Printf("  %-12s %v\n", n, err)
			continue
		}
		fmt.Printf("  %-12s path %.0f cells\n", n, m.TotalLen)
	}
	if m, err := game.LoadBoss(); err == nil {
		fmt.Printf("  %-12s path %.0f cells  (the heart: unseals after all floors, per renown)\n", "heart", m.TotalLen)
	}
	fmt.Println("procedural: -maze <seed>")
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "termtd: "+format+"\n", args...)
	os.Exit(1)
}
