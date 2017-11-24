package raweb

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ConfigService struct {
	config Config
}

func LoadConfigFromFile(config Config, fileName string) error {
	f, err := ioutil.ReadFile(fileName)
	if err == nil {
		return LoadConfigFromBytes(config, f)
	}
	return err
}

func LoadConfigFromBytes(config Config, bytes []byte) error {
	return yaml.Unmarshal(bytes, config)
}

func (service *ConfigService) Start(config Config) error {
	service.config = config
	return nil
}
