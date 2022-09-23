package library

import (
	"go.skymyer.dev/show-control/utils"
)

func FixtureLibraryFromConfig(file string) (Fixtures, error) {
	var fixtures = make(map[string]*Fixture)
	if err := utils.LoadFromFile(file, &fixtures); err != nil {
		return nil, err
	}

	// TODO: add validators:: value/ranges, head refs, features, ...

	return fixtures, nil
}

type Fixtures map[string]*Fixture

type Fixture struct {
	Manufacturer string
	Model        string
	Heads        map[string]*Head
	Modes        map[string]map[string]*Channel
}

type Head struct {
	Name string
}

type Channel struct {
	Channel     int
	Default     int
	Description string
	Functions   map[string]*ChannelFunction
}

type ChannelFunction struct {
	Features []string
	Min      int
	Max      int
}
