package manager

import (
	"context"
	"errors"
	"testing"

	"github.com/pydio/cells-runtime/config"
)

func resetIntegrations(t *testing.T) {
	integration.Lock()
	previousContributors := integration.templateData
	previousUpgrade := integration.upgrade
	integration.templateData = nil
	integration.upgrade = nil
	integration.Unlock()
	t.Cleanup(func() {
		integration.Lock()
		integration.templateData = previousContributors
		integration.upgrade = previousUpgrade
		integration.Unlock()
	})
}

func TestTemplateDataContributorsCompose(t *testing.T) {
	resetIntegrations(t)
	RegisterTemplateDataContributor(func(context.Context, string, config.Store) (map[string]any, error) {
		return map[string]any{"application": "example"}, nil
	})
	data := map[string]any{"runtime": true}
	if err := contributeTemplateData(context.Background(), "", nil, data); err != nil {
		t.Fatal(err)
	}
	if data["runtime"] != true || data["application"] != "example" {
		t.Fatalf("unexpected template data: %#v", data)
	}
}

func TestConfigUpgradeFailureIsReturned(t *testing.T) {
	resetIntegrations(t)
	want := errors.New("upgrade failed")
	RegisterConfigUpgrade(func(context.Context, config.Store) (bool, string, string, error) {
		return false, "", "", want
	})
	_, _, _, got := upgradeConfig(context.Background(), nil)
	if !errors.Is(got, want) {
		t.Fatalf("upgrade error = %v, want %v", got, want)
	}
}
