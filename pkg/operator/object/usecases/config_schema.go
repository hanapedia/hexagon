package usecases

import "github.com/hanapedia/the-bench/pkg/operator/constants"

type Route struct {
	Route  string               `json:"route,omitempty"`
	Method constants.HttpMethod `json:"method,omitempty"`
	Weight int32                `json:"weight,omitempty"`
}

type Config struct {
	Vus       int32  `json:"vus,omitempty"`
	Duration  string `json:"duration,omitempty"`
	UrlPrefix string `json:"urlPrefix,omitempty"`
}
