package config

import "testing"
import "fmt"

const configpath = "/tmp/default.yaml"

func TestWriteDefault(t *testing.T) {
	err := WriteConfig(DefaultConfig, configpath)
	if err != nil {
		t.Error("Writing default config errored\n", err)
	}
}

func TestReadDefault(t *testing.T) {
	TestWriteDefault(t)
	c, err := LoadConfig(configpath)
	if err != nil {
		t.Error("Reading default config errored\n", err)
	}

	fmt.Println(c)
}
