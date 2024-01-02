package config

type Setup struct {
	Name      string
	Devices   map[string]*Device
	Universes map[string]*Universe
	Fixtures  map[string]*Fixture
	Groups    map[string]*Group
	Artnet    Artnet
}

type Device struct {
	Driver string
	Device string
}

type Universe struct {
	Output Output
}

type Output struct {
	Device   string
	Universe int
}

type Fixture struct {
	Kind     string
	Mode     string
	Universe string
	Address  int

	// If mode == virtual, link it to the real device
	Real string
	// And supply a feature map
	Map map[string]string
}

type Group struct {
	Members []string
}

type Artnet struct {
	Enabled   bool
	Network string
}
