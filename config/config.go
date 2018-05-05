package config

import (
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

var CONFIG_LOCATION = func() string {
	// Maybe do something smarter here
	return path.Join(os.Getenv("HOME"), ".wmsrc.yaml")
}()

var lastConfig, _ = getEasyConfigInternal()
var lastRead = time.Now()
var mutex sync.Mutex

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

func getEasyConfigInternal() (Configuration, error) {
	conf, err := LoadConfig(CONFIG_LOCATION)
	if err == nil {
		return conf, err
	}
	err = WriteDefault(CONFIG_LOCATION)
	if err != nil {
		return conf, err
	}
	return LoadConfig(CONFIG_LOCATION)
}

func GetEasyConfig() (Configuration, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if time.Now().Sub(lastRead) > time.Minute*5 {
		conf, err := getEasyConfigInternal()
		lastConfig = conf
		lastRead = time.Now()
		return conf, err
	}
	return lastConfig, nil
}
