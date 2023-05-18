package storage

import "context"

type Storage interface {
	// Define the methods you need to interact with the storage backend
	CreateOrUpdateResource(ctx context.Context, resourceID string, data interface{}) error
	GetResource(ctx context.Context, resourceID string) (interface{}, error)
	DeleteResource(ctx context.Context, resourceID string) error
}
