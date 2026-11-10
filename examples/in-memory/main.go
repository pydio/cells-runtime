// Command in-memory demonstrates using cells-runtime without the Cells application.
package main

import (
	"context"
	"fmt"

	"github.com/spf13/viper"

	"github.com/pydio/cells-runtime/registry"
	"github.com/pydio/cells-runtime/runtime"
	"github.com/pydio/cells-runtime/runtime/manager"
	"github.com/pydio/cells-runtime/service"
	"github.com/pydio/cells-runtime/utils/propagator"

	_ "github.com/pydio/cells-runtime/config/memory"
	_ "github.com/pydio/cells-runtime/registry/config"
)

func main() {
	settings := viper.New()
	settings.Set(runtime.KeyKeyring, "mem://")
	settings.Set(runtime.KeyConfig, "mem://")
	settings.Set(runtime.KeyRegistry, "mem://")
	settings.Set(runtime.KeyBootstrapRoot, "")
	runtime.SetRuntime(settings)
	runtime.ConfigureApplication(runtime.ApplicationOptions{
		Vendor:           "Example",
		UnixVendor:       "example",
		Name:             "cells-runtime-example",
		ServiceNamespace: "example.grpc.",
		KeyringService:   "example.grpc.keyring",
		KeyringKey:       "master",
		LoggerNamespace:  "example",
		PackageLabel:     "Runtime Example",
		Version:          "0.1.0",
	})

	ctx := context.Background()
	mgr, err := manager.NewManager(ctx, "example")
	if err != nil {
		panic(err)
	}
	defer mgr.StopAll()

	ctx = propagator.With(ctx, registry.ContextKey, mgr.Registry())
	svc := service.NewService(service.Name("example.echo"), service.Context(ctx))
	if err := svc.Start(registry.WithContextR(ctx)); err != nil {
		panic(err)
	}
	fmt.Println("service started")
	if err := svc.Stop(registry.WithContextR(ctx)); err != nil {
		panic(err)
	}
	fmt.Println("service stopped")
}
