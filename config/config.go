package config

import "gopkg.in/yaml.v2"
import "io/ioutil"

// LoadConfig loads configuration at file in yaml
func LoadConfig(file string) (Configuration, error) {
	conf := Configuration{read: file}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal(data, &conf)
	// if err != nil {
	// 	return conf, err
	// }

	return conf, err
}

// WriteConfig writes configuration out into yaml file
func WriteConfig(c Configuration, file string) error {

	data, err := yaml.Marshal(c)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, data, 0644)
	return err
}

// WriteDefault writes out default configuration to file in yaml
func WriteDefault(file string) error {
	return WriteConfig(DefaultConfig, file)
}
