package model

type ConfigLoader interface {
	Load() (ServiceUnitConfig, error)
}
