package sql

import (
	"os"
	"time"

	"github.com/pydio/cells-runtime/runtime"
)

var (
	DefaultConnectionTimeout = 30 * time.Second
	LongConnectionTimeout    = 10 * time.Minute
)

func init() {
	runtime.RegisterApplicationEnvVariable("SQL_DEFAULT_CONN", "30s", "Default SQL connection timeout")
	runtime.RegisterApplicationEnvVariable("SQL_LONG_CONN", "10m", "Default SQL long connections timeout (for reading big tables)")
	runtime.RegisterApplicationConfigurator(func(runtime.ApplicationOptions) {
		if dc := os.Getenv(runtime.EnvironmentName("SQL_DEFAULT_CONN")); dc != "" {
			if ddc, e := time.ParseDuration(dc); e == nil {
				DefaultConnectionTimeout = ddc
			}
		}

		if dc := os.Getenv(runtime.EnvironmentName("SQL_LONG_CONN")); dc != "" {
			if ddc, e := time.ParseDuration(dc); e == nil {
				LongConnectionTimeout = ddc
			}
		}
	})
}
