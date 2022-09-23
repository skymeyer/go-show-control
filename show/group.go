package show

import "go.skymyer.dev/show-control/config"

func GetFixtureCollection(f map[string]*Fixture, g map[string]*config.Group,
	fixtures []string, groups []string) map[string]*Fixture {

	var list = make(map[string]*Fixture)

	// Groups
	for _, group := range groups {
		for _, i := range g[group].Members {
			list[i] = f[i]
		}
	}

	// Individual fixtures
	for _, i := range fixtures {
		list[i] = f[i]
	}

	// Flatten result
	return list
}
