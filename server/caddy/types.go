package caddy

import "github.com/rs/cors"

// ActiveRoute is the runtime representation of an HTTP route.
type ActiveRoute struct {
	Path         string
	HeaderMods   []any
	RewriteRules []any
	Upstreams    []any
	WebSocket    bool
	CorsOptions  any
}

// ActiveSite is the runtime representation consumed by the Caddy adapter.
// Applications convert their persisted site configuration into this type.
type ActiveSite struct {
	Binds                 []string
	Hash                  string
	Maintenance           bool
	MaintenanceConditions []string
	Routes                []*ActiveRoute
	CorsOptions           *cors.Options
	TLS                   string
	SSLRedirect           bool
	RedirectTargets       map[string]string
	MuxMode               bool
	Log                   string
	LogLevel              string
}

// Redirects returns the HTTP-to-HTTPS redirects for this site.
func (s *ActiveSite) Redirects() map[string]string {
	return s.RedirectTargets
}
