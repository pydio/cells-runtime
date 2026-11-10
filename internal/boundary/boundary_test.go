package boundary

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRuntimeDoesNotImportCells(t *testing.T) {
	root := filepath.Join("..", "..")
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" && filepath.Base(path) != "go.mod" {
			return nil
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		forbidden := strings.Join([]string{"github.com", "pydio", "cells", "v5"}, "/")
		if strings.Contains(string(contents), forbidden) {
			t.Errorf("%s imports the Cells application module", path)
		}
		productDefaults := []string{
			strings.Join([]string{"CELLS", "_"}, ""),
			strings.Join([]string{"PYDIO", "_"}, ""),
			strings.Join([]string{`"`, "pydio", "."}, ""),
		}
		for _, productDefault := range productDefaults {
			if strings.Contains(string(contents), productDefault) {
				t.Errorf("%s contains Cells application default %q", path, productDefault)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
