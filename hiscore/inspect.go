package hiscore

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// SaveSchema names the existing on-disk contract. Diagnostic and packaging
// metadata use it; it does not alter or migrate the player's save files.
const SaveSchema = "termtd-v1"

type SaveStatus struct {
	Kind   string `json:"kind"`
	Path   string `json:"path"`
	Source string `json:"source,omitempty"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}
type SaveReport struct {
	Schema  string       `json:"schema"`
	Files   []SaveStatus `json:"files"`
	Healthy bool         `json:"healthy"`
}

// Inspect reports compatibility with this build without creating, migrating,
// repairing or rewriting saves. Missing saves are a normal first-run state.
func Inspect() (SaveReport, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return SaveReport{}, err
	}
	report := SaveReport{Schema: SaveSchema, Healthy: true}
	for _, kind := range []string{"hiscores", "journal", "lair"} {
		path := filepath.Join(home, ".termtd-"+kind+".json")
		source := path
		st, err := os.Lstat(source)
		if os.IsNotExist(err) {
			source = legacyPath(path)
			st, err = os.Lstat(source)
		}
		item := SaveStatus{Kind: kind, Path: path}
		if os.IsNotExist(err) {
			item.Status = "missing"
			report.Files = append(report.Files, item)
			continue
		}
		item.Source = source
		if err != nil {
			item.Status = "unreadable"
			item.Detail = "cannot inspect save file"
		} else if !st.Mode().IsRegular() {
			item.Status = "unsupported"
			item.Detail = "save is not a regular file"
		} else if st.Size() > 32<<20 {
			item.Status = "unsupported"
			item.Detail = "save exceeds diagnostic size limit"
		} else {
			f, e := os.Open(source)
			if e != nil {
				item.Status = "unreadable"
				item.Detail = "cannot read save file"
			} else {
				data, e := io.ReadAll(io.LimitReader(f, 32<<20+1))
				f.Close()
				if e != nil || len(data) > 32<<20 {
					item.Status = "unreadable"
					item.Detail = "cannot read bounded save"
				} else {
					var value any
					switch kind {
					case "hiscores":
						value = &Table{}
					case "journal":
						value = &Journal{}
					case "lair":
						value = &Lair{}
					}
					if e = decodeSave(data, value); e != nil {
						item.Status = "invalid"
						item.Detail = "save cannot be decoded by this build"
					} else {
						item.Status = "ok"
						if source != path {
							item.Status = "legacy"
						}
					}
				}
			}
		}
		if item.Status != "ok" && item.Status != "legacy" {
			report.Healthy = false
		}
		report.Files = append(report.Files, item)
	}
	return report, nil
}
func WriteReport(w io.Writer, report SaveReport, asJSON bool) error {
	if asJSON {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(report)
	}
	fmt.Fprintln(w, "Save format:", report.Schema)
	for _, item := range report.Files {
		path := item.Source
		if path == "" {
			path = item.Path
		}
		fmt.Fprintf(w, "%s: %s — %s\n", item.Kind, item.Status, path)
		if item.Detail != "" {
			fmt.Fprintln(w, " ", item.Detail)
		}
	}
	_, e := fmt.Fprintln(w, "Read-only check. No saves were changed.")
	return e
}
