package manager

import (
	"context"
	"sync"

	"github.com/pydio/cells-runtime/config"
)

// TemplateDataContributor adds application-specific values to bootstrap templates.
type TemplateDataContributor func(context.Context, string, config.Store) (map[string]any, error)

// ConfigUpgrade upgrades an application configuration after the runtime store opens.
// The returned user and message are used when persisting an upgraded configuration.
type ConfigUpgrade func(context.Context, config.Store) (upgraded bool, user, message string, err error)

var integration = struct {
	sync.RWMutex
	templateData []TemplateDataContributor
	upgrade      ConfigUpgrade
}{}

// RegisterTemplateDataContributor registers application bootstrap data.
func RegisterTemplateDataContributor(contributor TemplateDataContributor) {
	integration.Lock()
	defer integration.Unlock()
	integration.templateData = append(integration.templateData, contributor)
}

// RegisterConfigUpgrade registers the application configuration upgrader.
func RegisterConfigUpgrade(upgrade ConfigUpgrade) {
	integration.Lock()
	defer integration.Unlock()
	integration.upgrade = upgrade
}

func contributeTemplateData(ctx context.Context, tpl string, store config.Store, data map[string]any) error {
	integration.RLock()
	contributors := append([]TemplateDataContributor(nil), integration.templateData...)
	integration.RUnlock()
	for _, contributor := range contributors {
		values, err := contributor(ctx, tpl, store)
		if err != nil {
			return err
		}
		for key, value := range values {
			data[key] = value
		}
	}
	return nil
}

func upgradeConfig(ctx context.Context, store config.Store) (bool, string, string, error) {
	integration.RLock()
	upgrade := integration.upgrade
	integration.RUnlock()
	if upgrade == nil {
		return false, "", "", nil
	}
	return upgrade(ctx, store)
}
