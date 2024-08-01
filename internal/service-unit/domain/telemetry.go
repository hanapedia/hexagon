package domain

import model "github.com/hanapedia/hexagon/pkg/api/v1"

type Status float64

const (
	Ok Status = iota
	Err
)

func (s Status) AsString() string {
	switch s {
	case Ok:
		return "ok"
	case Err:
		return "error"
	}
	return "incomplete"
}

type PrimaryLabels struct {
	ServiceName string
	Variant     string
	Action      string
	Route       string
	Topic       string
}

type SecondaryLabels struct {
	DstServiceName string
	DstVariant     string
	DstAction      string
	DstRoute       string
	DstTopic       string
}

type TelemetryContext struct {
	PrimaryLabels
	SecondaryLabels
}

var AllLabelKeys = []string{
	"service",
	"variant",
	"action",
	"route",
	"topic",
	"dst_service",
	"dst_variant",
	"dst_action",
	"dst_route",
	"dst_topic",
}

var PrimaryLabelKeys = []string{
	"service",
	"variant",
	"action",
	"route",
	"topic",
}

var SecondaryLabelKeys = []string{
	"dst_service",
	"dst_variant",
	"dst_action",
	"dst_route",
	"dst_topic",
}

func (tx TelemetryContext) AsMap() map[string]string {
	return map[string]string{
		"service":     tx.ServiceName,
		"variant":     tx.Variant,
		"action":      tx.Action,
		"route":       tx.Route,
		"topic":       tx.Topic,
		"dst_service": tx.DstServiceName,
		"dst_variant": tx.DstVariant,
		"dst_action":  tx.DstAction,
		"dst_route":   tx.DstRoute,
		"dst_topic":   tx.DstTopic,
	}
}

func (pl PrimaryLabels) GetPrimaryLabels() map[string]string {
	return map[string]string{
		"service": pl.ServiceName,
		"variant": pl.Variant,
		"action":  pl.Action,
		"route":   pl.Route,
		"topic":   pl.Topic,
	}
}

func (sl SecondaryLabels) GetSecondaryLabels() map[string]string {
	return map[string]string{
		"dst_service": sl.DstServiceName,
		"dst_variant": sl.DstVariant,
		"dst_action":  sl.DstAction,
		"dst_route":   sl.DstRoute,
		"dst_topic":   sl.DstTopic,
	}
}

func NewTelemetryContext(serviceName string, primaryConfig *model.PrimaryAdapterSpec, secondaryConfig *model.SecondaryAdapterConfig) TelemetryContext {
	telCtx := TelemetryContext{
		PrimaryLabels: PrimaryLabels{
			ServiceName: serviceName,
		},
	}

	if primaryConfig.ServerConfig != nil {
		telCtx.Variant = string(primaryConfig.ServerConfig.Variant)
		telCtx.Action = string(primaryConfig.ServerConfig.Action)
		telCtx.Route = primaryConfig.ServerConfig.Route
	} else if primaryConfig.ConsumerConfig != nil {
		telCtx.Variant = string(primaryConfig.ConsumerConfig.Variant)
		telCtx.Topic = primaryConfig.ConsumerConfig.Topic
	}

	if secondaryConfig.InvocationConfig != nil {
		telCtx.DstServiceName = secondaryConfig.InvocationConfig.Service
		telCtx.DstVariant = string(secondaryConfig.InvocationConfig.Variant)
		telCtx.DstAction = string(secondaryConfig.InvocationConfig.Action)
		telCtx.DstRoute = secondaryConfig.InvocationConfig.Route
	} else if secondaryConfig.ProducerConfig != nil {
		telCtx.DstVariant = string(secondaryConfig.ProducerConfig.Variant)
		telCtx.DstTopic = secondaryConfig.ProducerConfig.Topic
	} else if secondaryConfig.RepositoryConfig != nil {
		telCtx.DstServiceName = secondaryConfig.RepositoryConfig.Name
		telCtx.DstVariant = string(secondaryConfig.RepositoryConfig.Variant)
		telCtx.DstAction = string(secondaryConfig.RepositoryConfig.Action)
	} else if secondaryConfig.StressorConfig != nil {
		telCtx.DstServiceName = string(secondaryConfig.StressorConfig.Name)
		telCtx.DstVariant = string(secondaryConfig.StressorConfig.Variant)
	}

	return telCtx
}
