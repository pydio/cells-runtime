package prometheus

import (
	"context"
	"sync"
)

var siteURL = struct {
	sync.RWMutex
	provider func(context.Context) string
}{}

// RegisterExternalURLProvider installs the application's externally visible URL resolver.
func RegisterExternalURLProvider(provider func(context.Context) string) {
	siteURL.Lock()
	defer siteURL.Unlock()
	siteURL.provider = provider
}

func externalURL(ctx context.Context) string {
	siteURL.RLock()
	defer siteURL.RUnlock()
	if siteURL.provider == nil {
		return ""
	}
	return siteURL.provider(ctx)
}
