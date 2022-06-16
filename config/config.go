package config

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type RawConfig map[string]interface{}

func (p *RawConfig) UnpackRaw(to interface{}) error {

	mcfg := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc()),
		WeaklyTypedInput: true,
		Metadata:         nil,
		Result:           to,
	}

	d, err := mapstructure.NewDecoder(mcfg)
	if err != nil {
		return fmt.Errorf("error create new config decoder: %w", err)
	}

	if err := d.Decode(p); err != nil {
		return fmt.Errorf("error decoding config: %w", err)
	}
	return nil
}

func (p *RawConfig) UnpackAttribute(value string, to interface{}) error {

	mcfg := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc()),
		WeaklyTypedInput: true,
		Metadata:         nil,
		Result:           to,
	}

	d, err := mapstructure.NewDecoder(mcfg)
	if err != nil {
		return fmt.Errorf("error create new config decoder: %w", err)
	}

	c := *p

	if err := d.Decode(c[value]); err != nil {
		return fmt.Errorf("error decoding config: %w", err)
	}

	return nil
}

type ConfigNamespace struct {
	name   string
	config *RawConfig
}

func (c ConfigNamespace) Name() string {
	return c.name
}

func (c ConfigNamespace) Config() *RawConfig {
	return c.config
}

func UnpackNamespace(value string, cfg *RawConfig) (ConfigNamespace, error) {

	config := ConfigNamespace{}

	c := *cfg

	if _, ok := c[value]; !ok {
		return ConfigNamespace{}, fmt.Errorf("no %s configured in config", value)
	}

	output := c[value]

	for k, v := range output.(map[string]interface{}) {
		config.name = k

		if err := mapstructure.Decode(v, &config.config); err != nil {
			return ConfigNamespace{}, fmt.Errorf("unable to decode %s namespace: %w", value, err)
		}

	}

	return config, nil

}
