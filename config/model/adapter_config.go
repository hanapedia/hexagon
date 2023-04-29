package model

import "github.com/hanapedia/the-bench/config/constants"

// Config fields for stateful services
type StatelessAdapterConfig struct {
	Variant constants.StatelessAdapterVariant `yaml:"variant" validate:"required,oneof=rest grpc"`
	Service string                            `yaml:"service,omitempty" validate:"required"`
	Action  constants.Action                  `yaml:"action" validate:"required,oneof=read write"`
	Route   string                            `yaml:"route" validate:"required"`
}

// Config fields for stateful services
type StatefulAdapterConfig struct {
	Name    string                           `yaml:"name" validate:"required"`
	Variant constants.StatefulAdapterVariant `yaml:"variant" validate:"required,oneof=mongo postgre"`
	Action  constants.Action                 `yaml:"action" validate:"omitempty,oneof=read write"`
	Size    string                           `yaml:"size" validate:"omitempty,oneof=small medium large"`
}

// Config fields for Brokers
type BrokerAdapterConfig struct {
	Variant constants.BrokerAdapterVariant `yaml:"variant" validate:"required,oneof=kafka rabbitmq pulsar"`
	Topic   string                         `yaml:"topic" validate:"required"`
}

// Config fields for Internal services
type InternalAdapterConfig struct {
	Name     string `yaml:"name" validate:"required"`
	Resource string `yaml:"resource" validate:"required,oneof=cpu memory disk network"`
	Duration string `yaml:"duration" validate:"required,oneof=small medium large"`
	Load     string `yaml:"load" validate:"required,oneof=small medium large"`
}
