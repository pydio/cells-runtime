package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"

	"github.com/pydio/cells-runtime/utils/propagator"
)

// AuthModifierFactory supplies application authentication for public gRPC servers.
type AuthModifierFactory func(context.Context) propagator.IncomingContextModifier

// TLSConfigProvider supplies certificates for secure gRPC servers.
type TLSConfigProvider func(context.Context) (*tls.Config, error)

var integrations = struct {
	sync.RWMutex
	auth AuthModifierFactory
	tls  TLSConfigProvider
}{}

// RegisterAuthModifierFactory installs application authentication handling.
func RegisterAuthModifierFactory(factory AuthModifierFactory) {
	integrations.Lock()
	defer integrations.Unlock()
	integrations.auth = factory
}

// RegisterTLSConfigProvider installs application TLS configuration handling.
func RegisterTLSConfigProvider(provider TLSConfigProvider) {
	integrations.Lock()
	defer integrations.Unlock()
	integrations.tls = provider
}

func authModifier(ctx context.Context) propagator.IncomingContextModifier {
	integrations.RLock()
	factory := integrations.auth
	integrations.RUnlock()
	if factory != nil {
		return factory(ctx)
	}
	return func(in context.Context) (context.Context, bool, error) {
		return propagator.ForkContext(in, ctx), false, nil
	}
}

func tlsConfig(ctx context.Context) (*tls.Config, error) {
	integrations.RLock()
	provider := integrations.tls
	integrations.RUnlock()
	if provider == nil {
		return nil, fmt.Errorf("secure gRPC server requested without a TLS configuration provider")
	}
	return provider(ctx)
}
