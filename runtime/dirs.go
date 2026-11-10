/*
 * Copyright (c) 2019-2026. Abstrium SAS <team (at) pydio.com>
 * This file is part of Pydio Cells Runtime.
 *
 * Pydio Cells Runtime is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */

package runtime

import (
	"errors"
	"os"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"sync"

	hashiversion "github.com/hashicorp/go-version"
	"github.com/shibukawa/configdir"

	common "github.com/pydio/cells-runtime"
)

type ApplicationDirType int

const (
	ApplicationDirData ApplicationDirType = iota
	ApplicationDirLogs
	ApplicationDirServices
)

// ApplicationOptions contains application policy used to resolve runtime
// directories. The zero value intentionally contains no Cells-specific names.
type ApplicationOptions struct {
	Vendor      string
	UnixVendor  string
	Name        string
	WorkingDir  string
	DataDir     string
	LogDir      string
	ServicesDir string
	WorkingEnv  string
	DataEnv     string
	LogEnv      string
	ServicesEnv string
	EnvPrefix   string
	// ServiceNamespace prefixes infrastructure service names. KeyringService
	// and KeyringKey identify the runtime encryption key.
	ServiceNamespace     string
	WebNamespace         string
	RESTNamespace        string
	GenericNamespace     string
	RegistryTopic        string
	MetadataPrefix       string
	ContextUserKey       string
	SiteHashHeader       string
	ForwardedHeaders     []string
	SilentStreamServices []string
	ConfigFileName       string
	VaultFileName        string
	KeyringFileName      string
	CertStorePath        string
	KeyringService       string
	KeyringKey           string
	LoggerNamespace      string
	PackageLabel         string
	Version              string
	BuildStamp           string
	BuildRevision        string
}

var application = struct {
	sync.RWMutex
	options ApplicationOptions
}{options: ApplicationOptions{
	Vendor:           "Runtime",
	UnixVendor:       "runtime",
	Name:             "runtime",
	WorkingEnv:       "RUNTIME_WORKING_DIR",
	DataEnv:          "RUNTIME_DATA_DIR",
	LogEnv:           "RUNTIME_LOG_DIR",
	ServicesEnv:      "RUNTIME_SERVICES_DIR",
	EnvPrefix:        "RUNTIME",
	ServiceNamespace: "runtime.grpc.",
	KeyringService:   "runtime.grpc.user-key",
	KeyringKey:       "keyring.master",
	LoggerNamespace:  "runtime",
	PackageLabel:     "Cells Runtime",
	Version:          "0.0.0",
}}

var applicationConfigurers = struct {
	sync.RWMutex
	callbacks []func(ApplicationOptions)
}{}

// RegisterApplicationConfigurator registers a callback that applies
// application policy after ConfigureApplication is called.
func RegisterApplicationConfigurator(callback func(ApplicationOptions)) {
	applicationConfigurers.Lock()
	defer applicationConfigurers.Unlock()
	applicationConfigurers.callbacks = append(applicationConfigurers.callbacks, callback)
}

// ConfigureApplication replaces directory and environment policy. Applications
// should call it before constructing a manager.
func ConfigureApplication(options ApplicationOptions) {
	application.Lock()
	application.options = options
	application.Unlock()
	common.ConfigureRuntimeNaming(common.RuntimeNaming{
		GRPCNamespace:    options.ServiceNamespace,
		WebNamespace:     options.WebNamespace,
		RESTNamespace:    options.RESTNamespace,
		GenericNamespace: options.GenericNamespace,
		RegistryTopic:    options.RegistryTopic,
		MetadataPrefix:   options.MetadataPrefix,
		ContextUserKey:   options.ContextUserKey,
		SiteHashHeader:   options.SiteHashHeader,
		ForwardedHeaders: options.ForwardedHeaders,
	})
	applicationConfigurers.RLock()
	callbacks := append([]func(ApplicationOptions){}, applicationConfigurers.callbacks...)
	applicationConfigurers.RUnlock()
	for _, callback := range callbacks {
		callback(options)
	}
	if options.ConfigFileName != "" {
		DefaultConfigFileName = options.ConfigFileName
	}
	if options.VaultFileName != "" {
		DefaultVaultFileName = options.VaultFileName
	}
	if options.KeyringFileName != "" {
		DefaultKeyringFileName = options.KeyringFileName
	}
	if options.CertStorePath != "" {
		DefaultCertStorePath = options.CertStorePath
	}
	refreshDefaultPaths()
}

func applicationOptions() ApplicationOptions {
	application.RLock()
	defer application.RUnlock()
	return application.options
}

// Application returns a copy of the configured application policy.
func Application() ApplicationOptions { return applicationOptions() }

// EnvironmentName returns an application-scoped environment variable name.
func EnvironmentName(suffix string) string {
	prefix := strings.TrimSuffix(Application().EnvPrefix, "_")
	if prefix == "" {
		return suffix
	}
	return prefix + "_" + suffix
}

// ApplicationVersion returns the version supplied by the embedding application.
func ApplicationVersion() *hashiversion.Version {
	value := Application().Version
	parsed, err := hashiversion.NewVersion(value)
	if err != nil {
		parsed, _ = hashiversion.NewVersion("0.0.0")
	}
	return parsed
}

func configuredValue(explicit, environment string) string {
	if explicit != "" {
		return explicit
	}
	if environment != "" {
		return os.Getenv(environment)
	}
	return ""
}

// ApplicationWorkingDir resolves an application directory. It does not create
// directories and has no import-time filesystem side effects.
func ApplicationWorkingDir(dirType ...ApplicationDirType) string {
	o := applicationOptions()
	base := configuredValue(o.WorkingDir, o.WorkingEnv)
	if len(dirType) > 0 {
		var explicit, environment, suffix string
		switch dirType[0] {
		case ApplicationDirData:
			explicit, environment, suffix = o.DataDir, o.DataEnv, "data"
		case ApplicationDirLogs:
			explicit, environment, suffix = o.LogDir, o.LogEnv, "logs"
		case ApplicationDirServices:
			explicit, environment, suffix = o.ServicesDir, o.ServicesEnv, "services"
		}
		if value := configuredValue(explicit, environment); value != "" {
			return value
		}
		if base != "" {
			return filepath.Join(base, suffix)
		}
	}
	if base != "" {
		return base
	}
	vendor := o.Vendor
	if goruntime.GOOS == "linux" || goruntime.GOOS == "freebsd" {
		vendor = o.UnixVendor
	}
	if vendor == "" || o.Name == "" {
		return ""
	}
	configDirs := configdir.New(vendor, o.Name)
	folders := configDirs.QueryFolders(configdir.Global)
	if len(folders) == 0 {
		folders = configDirs.QueryFolders(configdir.Local)
	}
	if len(folders) == 0 {
		return ""
	}
	if len(dirType) == 0 {
		return folders[0].Path
	}
	suffix := map[ApplicationDirType]string{
		ApplicationDirData: "data", ApplicationDirLogs: "logs", ApplicationDirServices: "services",
	}[dirType[0]]
	return filepath.Join(folders[0].Path, suffix)
}

// ServiceDataDir creates and returns the directory for a service.
func ServiceDataDir(serviceName string) (string, error) {
	base := ApplicationWorkingDir(ApplicationDirServices)
	if base == "" {
		return "", errors.New("runtime application directories are not configured")
	}
	dir := filepath.Join(base, serviceName)
	return dir, os.MkdirAll(dir, 0o755)
}

// MustServiceDataDir creates and returns the directory for a service.
func MustServiceDataDir(serviceName string) string {
	dir, err := ServiceDataDir(serviceName)
	if err != nil {
		panic("cannot create service data dir: " + err.Error())
	}
	return dir
}
