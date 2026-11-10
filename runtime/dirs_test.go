package runtime

import (
	"path/filepath"
	"testing"
)

func TestApplicationDirectoriesAreConfiguredByApplication(t *testing.T) {
	previous := applicationOptions()
	t.Cleanup(func() { ConfigureApplication(previous) })
	ConfigureApplication(ApplicationOptions{WorkingDir: t.TempDir()})
	if got, want := ApplicationWorkingDir(ApplicationDirData), filepath.Join(ApplicationWorkingDir(), "data"); got != want {
		t.Fatalf("data directory = %q, want %q", got, want)
	}
}
