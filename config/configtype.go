package config

import "gopkg.in/yaml.v2"
import "fmt"

// Configuration for WMS
type Configuration struct {
	Title         string       `yaml:"title"`
	ExampleCities []string     `yaml:"example"`
	ExampleModi   []string     `yaml:"modus,flow"`
	Rendering     RenderConfig `yaml:"render,omitempty"`
	DTageLinks    []string     `yaml:"dtage"`

	read string
}

// RenderConfig sets up cities to render
type RenderConfig struct {
	// Cities to render
	Cities []string `yaml:"cites"`
	// Interval in number of 10 minutes
	Interval int `yaml:"interval"`
}

func (c Configuration) String() string {
	data, err := yaml.Marshal(&c)
	if err != nil {
		return "[error reading configuration]"
	}
	return fmt.Sprintf("%s\n{last read: '%s'}", string(data), c.read)
}
