package initialization_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hanapedia/hexagon/internal/service-unit/application/core/initialization"
	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	kafkaPrimary "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary/consumer/kafka"
	grpcPrimary "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary/server/grpc"
	restPrimary "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/primary/server/rest"
	grpcSecondary "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/invocation/grpc"
	restSecondary "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/invocation/rest"
	kafkaSecondary "github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/producer/kafka"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/repository/mongo"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/repository/redis"
	"github.com/hanapedia/hexagon/test/mock"
)

type TestCases struct {
	input string
	expected initialization.ServiceUnit
}

// TestServiceUnitSetup should test mapping between primary adapter and secondary adapter. It should
// 1. check that the correct primary adapters are setup
// 2. check that correct adapters are mapped to the primary adapter tasks
// 3. check that correct secondary adapter clients are setup
func TestServiceUnitSetup(t *testing.T) {
	// testCases compare that
	testCases := newTestCases()
	for _, tc := range testCases {
		config, err := mock.NewConfigLoader(tc.input).Load()
		if err != nil {
			t.Logf("Failed to load config from raw data: %s", err)
			t.Fail()
		}

		serviceUnit := initialization.NewServiceUnit(config)
		serviceUnit.Setup()
		serverAdapters := extractKeys(serviceUnit.ServerAdapters)
		consumerAdapters := extractKeys(serviceUnit.ConsumerAdapters)
		secondaryAdapterClients := extractKeys(serviceUnit.SecondaryAdapterClients)
		expectedServerAdapters := extractKeys(tc.expected.ServerAdapters)
		expectedConsumerAdapters := extractKeys(tc.expected.ConsumerAdapters)
		expectedSecondaryAdapterClients := extractKeys(tc.expected.SecondaryAdapterClients)
		assert.ElementsMatchf(t,
			expectedServerAdapters, serverAdapters,
			"Server adapters do not match. expected: %s, but got: %s",
			expectedServerAdapters, serverAdapters,
		)
		assert.ElementsMatchf(t,
			expectedConsumerAdapters, consumerAdapters,
			"Consumer adapters do not match. expected: %s, but got: %s",
			expectedConsumerAdapters, consumerAdapters,
		)
		assert.ElementsMatchf(t,
			expectedSecondaryAdapterClients, secondaryAdapterClients,
			"Secondary clients do not match. expected: %s, but got: %s",
			expectedSecondaryAdapterClients, secondaryAdapterClients,
		)
	}
}

func extractKeys(m interface{}) []string {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		return nil // Optionally, you could also return an error if you prefer
	}

	keysSlice := make([]string, 0, v.Len())
	for _, key := range v.MapKeys() {
		if key.Kind() == reflect.String {
			keysSlice = append(keysSlice, key.String())
		}
	}

	return keysSlice
}

func newTestCases() []TestCases {
	return []TestCases{
		{
			input:`version: dev
name: test
adapters:
- server:
    action: read
    variant: rest
    route: get
    payload:
      variant: large
  steps:
  - adapter:
      invocation:
        variant: rest
        service: chain1
        action: read
        route: get`,
			expected: initialization.ServiceUnit{
				Name: "test",
				ServerAdapters: map[string]ports.PrimaryPort{
					"rest": &restPrimary.RestServerAdapter{},
				},
				SecondaryAdapterClients: map[string]ports.SecondaryAdapterClient{
					"rest": &restSecondary.RestClient{},
				},
			},
		},
		{
			input: `version: dev
name: test
adapters:
- server:
    action: write
    variant: rest
    route: post
  steps:
  - adapter:
      invocation:
        variant: rest
        service: chain1
        action: write
        route: post
        payload:
          variant: large`,
			expected: initialization.ServiceUnit{
				Name: "test",
				ServerAdapters: map[string]ports.PrimaryPort{
					"rest": &restPrimary.RestServerAdapter{},
				},
				SecondaryAdapterClients: map[string]ports.SecondaryAdapterClient{
					"rest": &restSecondary.RestClient{},
				},
			},
		},
		{
			input: `version: dev
name: test
adapters:
- server:
    action: simpleRpc
    variant: grpc
    route: get
  steps:
  - adapter:
      invocation:
        variant: grpc
        service: grpcserver
        action: simpleRpc
        route: get`,
			expected: initialization.ServiceUnit{
				Name: "test",
				ServerAdapters: map[string]ports.PrimaryPort{
					"grpc": &grpcPrimary.GrpcServerAdapter{},
				},
				SecondaryAdapterClients: map[string]ports.SecondaryAdapterClient{
					"grpcserver.grpc": &grpcSecondary.GrpcClient{},
				},
			},
		},
		{
			input: `version: dev
name: test
adapters:
- server:
    action: clientStream
    variant: grpc
    route: get
  steps:
  - adapter:
      invocation:
        variant: grpc
        service: grpcserver
        action: clientStream
        route: get`,
			expected: initialization.ServiceUnit{
				Name: "test",
				ServerAdapters: map[string]ports.PrimaryPort{
					"grpc": &grpcPrimary.GrpcServerAdapter{},
				},
				SecondaryAdapterClients: map[string]ports.SecondaryAdapterClient{
					"grpcserver.grpc": &grpcSecondary.GrpcClient{},
				},
			},
		},
		{
			input: `version: dev
name: test
adapters:
- server:
    action: serverStream
    variant: grpc
    route: get
  steps:
  - adapter:
      invocation:
        variant: grpc
        service: grpcserver
        action: serverStream
        route: get`,
			expected: initialization.ServiceUnit{
				Name: "test",
				ServerAdapters: map[string]ports.PrimaryPort{
					"grpc": &grpcPrimary.GrpcServerAdapter{},
				},
				SecondaryAdapterClients: map[string]ports.SecondaryAdapterClient{
					"grpcserver.grpc": &grpcSecondary.GrpcClient{},
				},
			},
		},
		{
			input: `version: dev
name: test
adapters:
- server:
    action: biStream
    variant: grpc
    route: get
  steps:
  - adapter:
      invocation:
        variant: grpc
        service: grpcserver
        action: biStream
        route: get`,
			expected: initialization.ServiceUnit{
				Name: "test",
				ServerAdapters: map[string]ports.PrimaryPort{
					"grpc": &grpcPrimary.GrpcServerAdapter{},
				},
				SecondaryAdapterClients: map[string]ports.SecondaryAdapterClient{
					"grpcserver.grpc": &grpcSecondary.GrpcClient{},
				},
			},
		},
		{
			input: `version: dev
name: test
adapters:
- consumer:
    variant: kafka
    topic: topic1
  steps:
  - adapter:
      producer:
        variant: kafka
        topic: topic1`,
			expected: initialization.ServiceUnit{
				Name: "test",
				ConsumerAdapters: map[string]ports.PrimaryPort{
					"kafka.topic1": &kafkaPrimary.KafkaConsumerAdapter{},
				},
				SecondaryAdapterClients: map[string]ports.SecondaryAdapterClient{
					"kafka.topic1": &kafkaSecondary.KafkaProducerClient{},
				},
			},
		},
		{
			input: `version: dev
name: test
adapters:
- server:
    action: read
    variant: rest
    route: redis
  steps:
  - adapter:
      repository:
        name: redisrepo
        variant: redis
        action: read
        payload:
          variant: large
  - adapter:
      repository:
        name: redisrepo
        variant: redis
        action: write
        payload:
          variant: large
  - adapter:
      repository:
        name: mongorepo
        variant: mongo
        action: read
        payload:
          variant: large
  - adapter:
      repository:
        name: mongorepo
        variant: mongo
        action: write
        payload:
          variant: large`,
			expected: initialization.ServiceUnit{
				Name: "test",
				ServerAdapters: map[string]ports.PrimaryPort{
					"rest": &restPrimary.RestServerAdapter{},
				},
				SecondaryAdapterClients: map[string]ports.SecondaryAdapterClient{
					"redis.redisrepo": &redis.RedisClient{},
					"mongo.mongorepo": &mongo.MongoClient{},
				},
			},
		},
	}
}
