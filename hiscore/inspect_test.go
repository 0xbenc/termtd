package hiscore

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestInspectIsReadOnlyAndUsesCurrentThenLegacy(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	report, e := Inspect()
	if e != nil || !report.Healthy || len(report.Files) != 3 {
		t.Fatal(report, e)
	}
	for _, f := range report.Files {
		if f.Status != "missing" {
			t.Fatal(f)
		}
	}
	inputs := map[string][]byte{".termtd-hiscores.json": []byte(`{"rotunda":300}`), ".tdef-journal.json": []byte(`{"towers":{"archer":true}}`), ".termtd-lair.json": []byte(`{"tokens":8,"floors":{"rotunda:0":{"cleared":true}}}`)}
	for name, b := range inputs {
		os.WriteFile(filepath.Join(home, name), b, 0600)
	}
	report, e = Inspect()
	if e != nil || !report.Healthy || report.Files[1].Status != "legacy" {
		t.Fatal(report, e)
	}
	for name, b := range inputs {
		got, _ := os.ReadFile(filepath.Join(home, name))
		if !bytes.Equal(got, b) {
			t.Fatal("inspection rewrote save")
		}
	}
	os.WriteFile(filepath.Join(home, ".termtd-journal.json"), []byte(`broken`), 0600)
	report, e = Inspect()
	if e != nil || report.Healthy || report.Files[1].Status != "invalid" {
		t.Fatal("invalid current save fell back to legacy", report, e)
	}
	files, _ := os.ReadDir(home)
	if len(files) != 4 {
		t.Fatal("inspection created files", files)
	}
}
