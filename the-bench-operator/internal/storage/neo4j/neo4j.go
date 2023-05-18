package neo4j

import (
	"context"
	"github.com/hanapedia/the-bench/the-bench-operator/internal/storage"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jStorage struct {
	Driver neo4j.Driver
}

func New(driver neo4j.Driver) storage.Storage {
	return &Neo4jStorage{Driver: driver}
}

func NewNeo4jStorage(uri, username, password string) (storage.Storage, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}

	return &Neo4jStorage{Driver: driver}, nil
}

func (n *Neo4jStorage) CreateOrUpdateResource(ctx context.Context, resourceID string, data interface{}) error {
	// Implement the logic to create or update a resource in the Neo4j database
}

func (n *Neo4jStorage) GetResource(ctx context.Context, resourceID string) (interface{}, error) {
	// Implement the logic to get a resource from the Neo4j database
}

func (n *Neo4jStorage) DeleteResource(ctx context.Context, resourceID string) error {
	// Implement the logic to delete a resource from the Neo4j database
}
