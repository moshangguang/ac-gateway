package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type Config struct {
	AdminConfig *AdminConfig `yaml:"admin"`
	GateConfig  *GateConfig  `yaml:"gate"`
	EtcdConfig  *EtcdConfig  `yaml:"etcd"`
}

type GateConfig struct {
	Port    int `yaml:"port"`
	Timeout int `yaml:"timeout"`
}
type AdminConfig struct {
	Port int `yaml:"port"`
}

type EtcdConfig struct {
	Endpoints   []string `yaml:"endpoints"`
	DialTimeout int      `yaml:"dial-timeout"`
}

func NewConfig() Config {
	f, err := os.Open("config.yml")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()

	buf, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	conf := Config{}
	if err = yaml.Unmarshal(buf, &conf); err != nil {
		panic(err)
	}
	return conf
}
