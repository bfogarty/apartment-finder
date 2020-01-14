package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config stores the application's configuration.
type Config struct {
	Twilio struct {
		AccountSID string `yaml:"accountSid"`
		AuthToken  string `yaml:"authToken"`
		PhoneFrom  string `yaml:"phoneFrom"`
	} `yaml:"twilio"`
	Bitly struct {
		AccessToken string `yaml:"accessToken"`
	} `yaml:"bitly"`
	Notifications struct {
		RecipientPhones []string `yaml:"recipientPhones"`
	} `yaml:"notifications"`
}

// NewConfigFromYAML initializes a new config from a YAML file.
func NewConfigFromYAML(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, ErrLoadingConfig
	}

	defer f.Close()

	c := &Config{}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(c)
	if err != nil {
		return nil, ErrLoadingConfig
	}

	return c, nil
}
