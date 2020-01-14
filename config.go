package main

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Config stores the application's configuration.
type Config struct {
	Twilio struct {
		AccountSID string `yaml:"accountSid" envconfig:"TWILIO_ACCOUNT_SID"`
		AuthToken  string `yaml:"authToken" envconfig:"TWILIO_AUTH_TOKEN"`
		PhoneFrom  string `yaml:"phoneFrom" envconfig:"TWILIO_PHONE_FROM"`
	} `yaml:"twilio"`
	Bitly struct {
		AccessToken string `yaml:"accessToken" envconfig:"BITLY_ACCESS_TOKEN"`
	} `yaml:"bitly"`
	Notifications struct {
		RecipientPhones []string `yaml:"recipientPhones" envconfig:"RECIPIENT_PHONES"`
	} `yaml:"notifications"`
}

func loadConfigFromYAML(c *Config, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return ErrLoadingConfig
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(c)
	if err != nil {
		return err
	}

	return nil
}

func loadConfigFromEnv(c *Config) error {
	err := envconfig.Process("", c)
	if err != nil {
		return err
	}

	return nil
}

// NewConfig initializes a new config from the given YAML file.
// Environment variables can override the values in the file.
func NewConfig(path string) (*Config, error) {
	c := &Config{}

	err := loadConfigFromYAML(c, path)
	err = loadConfigFromEnv(c)
	if err != nil {
		return nil, ErrLoadingConfig
	}

	return c, nil
}
