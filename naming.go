package common

import "go.uber.org/zap/zapcore"

type (
	ServiceType string
	ServiceTag  string
	ServiceName string
)

const (
	ServiceLog      = "log"
	ServiceConfig   = "config"
	ServiceBroker   = "broker"
	ServiceRegistry = "registry"
	ServicePprof    = "pprof"
	ServiceOAuth    = "oauth"

	ServiceTagGateway = "gateway"

	CacheTypeShared     = "shared"
	CacheTypeLocal      = "local"
	QueueTypePersistent = "persistent"
	QueueTypeDebouncer  = "debouncer"

	CtxTargetServiceName  = "service"
	CtxGrpcClientCaller   = "grpc-client-caller"
	CtxGrpcSilentNotFound = "grpc-silent-not-found"
)

var (
	ServiceGrpcNamespace_    = "runtime.grpc."
	ServiceWebNamespace_     = "runtime.web."
	ServiceRestNamespace_    = "runtime.rest."
	ServiceGenericNamespace_ = "runtime.generic."

	ServiceRegistryGRPC = ServiceGrpcNamespace_ + ServiceRegistry
	ServiceBrokerGRPC   = ServiceGrpcNamespace_ + ServiceBroker
	ServiceConfigGRPC   = ServiceGrpcNamespace_ + ServiceConfig
	ServiceLogGRPC      = ServiceGrpcNamespace_ + ServiceLog

	TopicRegistryCommand = "runtime.registry.command"
	CtxCellsMetaPrefix   = "x-runtime-"
	PydioContextUserKey  = "runtime-user"
	XPydioSiteHash       = "X-Runtime-Site-Hash"
	XSpecialPydioHeaders []string
)

// RuntimeNaming contains application-owned protocol names used by generic
// infrastructure. Empty fields keep their current values.
type RuntimeNaming struct {
	GRPCNamespace    string
	WebNamespace     string
	RESTNamespace    string
	GenericNamespace string
	RegistryTopic    string
	MetadataPrefix   string
	ContextUserKey   string
	SiteHashHeader   string
	ForwardedHeaders []string
}

// ConfigureRuntimeNaming applies application protocol naming before runtime
// components are constructed.
func ConfigureRuntimeNaming(n RuntimeNaming) {
	if n.GRPCNamespace != "" {
		ServiceGrpcNamespace_ = n.GRPCNamespace
	}
	if n.WebNamespace != "" {
		ServiceWebNamespace_ = n.WebNamespace
	}
	if n.RESTNamespace != "" {
		ServiceRestNamespace_ = n.RESTNamespace
	}
	if n.GenericNamespace != "" {
		ServiceGenericNamespace_ = n.GenericNamespace
	}
	if n.RegistryTopic != "" {
		TopicRegistryCommand = n.RegistryTopic
	}
	if n.MetadataPrefix != "" {
		CtxCellsMetaPrefix = n.MetadataPrefix
	}
	if n.ContextUserKey != "" {
		PydioContextUserKey = n.ContextUserKey
	}
	if n.SiteHashHeader != "" {
		XPydioSiteHash = n.SiteHashHeader
	}
	if n.ForwardedHeaders != nil {
		XSpecialPydioHeaders = append([]string(nil), n.ForwardedHeaders...)
	}
	ServiceRegistryGRPC = ServiceGrpcNamespace_ + ServiceRegistry
	ServiceBrokerGRPC = ServiceGrpcNamespace_ + ServiceBroker
	ServiceConfigGRPC = ServiceGrpcNamespace_ + ServiceConfig
	ServiceLogGRPC = ServiceGrpcNamespace_ + ServiceLog
}

var (
	LogLevel  zapcore.Level
	LogJSON   bool
	LogToFile bool
)
