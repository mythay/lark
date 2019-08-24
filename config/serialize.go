package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func LoadCfgModbus(path string) (*CfgModbus, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var config CfgModbus
	if err := yaml.UnmarshalStrict(buf, &config); err != nil {
		return nil, err
	}
	return &config, nil

}
