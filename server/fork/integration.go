package fork

import "sync"

var customFlags = struct {
	sync.RWMutex
	provider func() []string
}{}

// RegisterCustomFlagProvider installs application flags forwarded to child processes.
func RegisterCustomFlagProvider(provider func() []string) {
	customFlags.Lock()
	defer customFlags.Unlock()
	customFlags.provider = provider
}

func forwardedFlags() []string {
	customFlags.RLock()
	defer customFlags.RUnlock()
	if customFlags.provider == nil {
		return nil
	}
	return customFlags.provider()
}
