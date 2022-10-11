package config

import (
	"fmt"
	"io/ioutil"

	"github.com/beeedge/beethings/cmd/rabbitmq-mqtt-plugin/app/options"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)

type Config struct {
	// TODO
	RabbitMQ *RabbitMQConfig `yaml:"rabbitMQ,omitempty"`
	Mqtt     *MqttConfig     `yaml:"mqtt,omitempty"`
}

type MqttConfig struct {
	Addr string `yaml:"addr,omitempty"`
	Port int    `yaml:"port,omitempty"`
	Qos  int    `yaml:"qos,omitempty"`
}

type RabbitMQConfig struct {
	Username     string `yaml:"username,omitempty"`
	Password     string `yaml:"password,omitempty"`
	Addr         string `yaml:"addr,omitempty"`
	Port         int    `yaml:"port,omitempty"`
	Exchange     string `yaml:"exchange,omitempty"`
	ExchangeType string `yaml:"exchangeType,omitempty"`
	RoutingKey   string `yaml:"routingKey,omitempty"`
	Reliable     bool   `yaml:"reliable,omitempty"`
	Queue        string `yaml:"queue,omitempty"`
	ConsumerTag  string `yaml:"consumerTag,omitempty"`
}

// loadConfig parses configuration file and returns
// an initialized Settings object and an error object if any. For instance if it
// cannot find the configuration file it will set the returned error appropriately.
func loadConfig(path string) (*Config, error) {
	c := &Config{}
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read configuration file %s, error %s", path, err)
	}
	if err = yaml.Unmarshal(contents, c); err != nil {
		return nil, fmt.Errorf("Failed to parse configuration, error %s", err)
	}
	if err = validate(c); err != nil {
		return nil, fmt.Errorf("Invalid configuration, error %s", err)
	}
	return c, nil
}

// validate the configuration
func validate(c *Config) error {
	// TODO: other validation check
	return nil
}

func NewRabbitmqMqttPluginConfig(o *options.Options) (*Config, error) {
	cfg, err := loadConfig(o.BaseCfgPath)
	if err != nil {
		klog.Errorf("loadConfig error %s", err.Error())
		return nil, err
	}
	return cfg, nil
}
