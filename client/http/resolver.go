package http

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/pydio/cells-runtime/client"
	"github.com/pydio/cells-runtime/config/routing"
	"github.com/pydio/cells-runtime/errors"
	"github.com/pydio/cells-runtime/middleware"
	"github.com/pydio/cells-runtime/registry"
	"github.com/pydio/cells-runtime/runtime"
	"github.com/pydio/cells-runtime/telemetry/log"
	"github.com/pydio/cells-runtime/utils/propagator"
)

var grpcTransport = &http.Transport{
	TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
	ForceAttemptHTTP2: true,
}

type Resolver interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) (bool, error)
	Init(ctx context.Context, serverID string, rr routing.RouteRegistrar)
	Stop()
}

// Mux wraps a standard MUX as an interface
type Mux interface {
	Handler(r *http.Request) (h http.Handler, pattern string)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// NewResolver creates an http resolver; If rootOnNotFound is set, non-matching patterns
// will be rewritten to base path and resolver will be run against it
type ResolverOptions struct {
	CoreServiceName    string
	GatewayGRPCService string
	Maintenance        http.Handler
}

func NewResolver(options ...ResolverOptions) Resolver {
	var resolved ResolverOptions
	if len(options) > 0 {
		resolved = options[0]
	}
	if resolved.CoreServiceName == "" {
		resolved.CoreServiceName = "runtime.web.mux"
	}
	return &resolver{options: resolved}
}

type cachedProxy struct {
	target string
	*httputil.ReverseProxy
}

type resolver struct {
	options      ResolverOptions
	c            grpc.ClientConnInterface
	r            registry.Registry
	coreCtx      context.Context
	rr           routing.RouteRegistrar
	s            Mux
	b            Balancer
	rc           client.ResolverCallback
	proxiesCache []*cachedProxy
}

func (m *resolver) Init(ctx context.Context, serverID string, rr routing.RouteRegistrar) {

	conn := runtime.GetClientConn(ctx)

	bal := NewBalancer(ctx, serverID)
	m.coreCtx = runtime.WithServiceName(runtime.AsCoreContext(ctx), m.options.CoreServiceName)
	m.c = conn
	m.rr = rr

	// Create balancer and listen to registry changes to build proxies on-demand
	m.b = bal
	var reg registry.Registry
	propagator.Get(ctx, registry.ContextKey, &reg)
	rc, _ := client.NewResolverCallback(reg)
	rc.Add(m.b.Build)

}

func (m *resolver) Stop() {}

func (m *resolver) ServeHTTP(w http.ResponseWriter, r *http.Request) (bool, error) {

	r, er := middleware.ApplyHTTPIncomingContextModifiers(r)
	if er != nil {
		return false, er
	}

	// This may rewrite the request, by resolving the underlying endpoints
	m.rr.ApplyRewrites(r)
	ctx := r.Context()

	// Special case for application/grpc
	if strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		if m.options.GatewayGRPCService == "" {
			return false, nil
		}
		proxyTarget, e := m.b.PickService(m.options.GatewayGRPCService)
		if e != nil {
			http.NotFound(w, r)
			return false, errors.New("cannot find grpc gateway")
		}
		proxy := m.getProxy(proxyTarget)
		// We assume that internally, the GRPCs service is serving self-signed
		proxy.Transport = grpcTransport
		// Wrap context and server request
		ctx = runtime.WithClientConn(ctx, m.c)
		proxy.ServeHTTP(w, r.WithContext(ctx))
		return true, nil
	}

	if m.options.Maintenance != nil && r.RequestURI == "/maintenance.html" && r.Header.Get("X-Maintenance-Redirect") != "" {
		m.options.Maintenance.ServeHTTP(w, r)
		return true, nil
	}

	ctx = runtime.WithClientConn(ctx, m.c)

	if m.s == nil {
		mu := http.NewServeMux()
		m.rr.IteratePatterns(m.coreCtx, func(pattern string, handler http.Handler) {
			// Wrap handler in OTEL middleware
			otl := middleware.HttpTracingMiddleware(pattern)
			mu.Handle(pattern, otl(handler))
		})
		m.s = mu
	}
	// fmt.Println(r.URL.Path, r.RequestURI)
	// Try to match pattern - Custom check for "/" : must be exactly /
	if h, pattern := m.s.Handler(r); len(pattern) > 0 && (pattern != "/" || r.URL.Path == "/") {
		h.ServeHTTP(w, r.WithContext(ctx))
		return true, nil
	}

	proxyTarget, e := m.b.PickEndpoint(r.URL.Path)
	if e == nil {
		m.getProxy(proxyTarget).ServeHTTP(w, r.WithContext(ctx))
		return true, nil
	}

	if m.rr.CanRewriteAndCatchAll(r) {
		if h, pat := m.s.Handler(r); len(pat) > 0 {
			h.ServeHTTP(w, r)
			return true, nil
		}
	}

	return false, nil
}

func (m *resolver) getProxy(url *url.URL) *httputil.ReverseProxy {
	for _, p := range m.proxiesCache {
		if p.target == url.String() {
			return p.ReverseProxy
		}
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		if err.Error() == "context canceled" {
			return
		}
		log.Logger(request.Context()).Error("Proxy Error :"+err.Error(), zap.Error(err))
		writer.WriteHeader(http.StatusBadGateway)
	}
	m.proxiesCache = append(m.proxiesCache, &cachedProxy{target: url.String(), ReverseProxy: proxy})
	return proxy
}

func (m *resolver) rewriteToRoot(r *http.Request) {
	r.URL.Path = "/"
	r.URL.RawPath = "/"
	r.RequestURI = r.URL.RequestURI()
}
