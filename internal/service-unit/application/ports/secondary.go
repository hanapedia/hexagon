package ports

import (
	"context"
)

// SecodaryPort provides common interface for all the external service resouce.
// Example resources include:
// - REST API routes
// - gRPC methods
// - Kafka topic
// - Database table
//
// It is intended to represent the individual endpoints on each exteranl service,
// not the services themselves; hence the name `SecodaryPort`
type SecodaryPort interface {
	// Call(context.Context) (string, error)
	Call(context.Context) SecondaryPortCallResult
	SetDestId(string)
	GetDestId() string
}

// Used to reuse clients to other serivces
// Wrapper interface, so the struct to implement this should have pointer to actual client
// types to implement this interface should have some sort of pointer to clients
type SecondaryAdapter interface {
	Close()
}
//
// type SecondaryPortError struct {
// 	SecondaryPort *SecodaryPort
// 	Error         error
// }

type SecondaryPortCallResult struct {
	Payload *string
	Error error
}

type SecondaryPortBase struct {
	destId string
}

func (spb *SecondaryPortBase) GetDestId() string {
	return spb.destId
}

func (spb *SecondaryPortBase) SetDestId(id string) {
	spb.destId = id
}
