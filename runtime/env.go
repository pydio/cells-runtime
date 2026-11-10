package runtime

import (
	"cmp"
	"fmt"
	"slices"
)

type envVar struct {
	key, description, def string
	isBool                bool
	applicationScoped     bool
}

// RegisterApplicationEnvVariable declares an environment variable whose name
// is prefixed by the embedding application's EnvPrefix.
func RegisterApplicationEnvVariable(suffix, def, description string, isBool ...bool) {
	RegisterEnvVariable(suffix, def, description, isBool...)
	ee[len(ee)-1].applicationScoped = true
}

var (
	ee []envVar
)

// RegisterEnvVariable can be called by any contributor to declare environment variables for documentation purpose
func RegisterEnvVariable(key, def, description string, isBool ...bool) {
	ib := false
	if len(isBool) > 0 {
		ib = isBool[0]
	}
	ee = append(ee, envVar{
		key:         key,
		def:         def,
		description: description,
		isBool:      ib,
	})
}

func DocRegisteredEnvVariables(excludes ...string) string {
	slices.SortStableFunc(ee, func(a, b envVar) int {
		return cmp.Compare(a.key, b.key)
	})
	s := ""
	for _, v := range ee {
		key := v.key
		if v.applicationScoped {
			key = EnvironmentName(key)
		}
		if slices.Contains(excludes, key) {
			continue
		}
		var def string
		if v.def != "" {
			def = " (" + v.def + ")"
		} else if v.isBool {
			def = " (false)"
		}
		s += fmt.Sprintf("  - %s%s: %s\n", key, def, v.description)
	}
	return s
}
