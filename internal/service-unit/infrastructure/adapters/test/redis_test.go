package test
/**/
/* import ( */
/* 	"context" */
/* 	"fmt" */
/* 	"strconv" */
/* 	"testing" */
/* 	"time" */
/**/
/* 	"github.com/docker/go-connections/nat" */
/* 	goRedis "github.com/redis/go-redis/v9" */
/* 	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/repository/redis" */
/* 	"github.com/hanapedia/hexagon/pkg/api/defaults" */
/* 	model "github.com/hanapedia/hexagon/pkg/api/v1" */
/* 	"github.com/hanapedia/hexagon/pkg/operator/logger" */
/* 	"github.com/testcontainers/testcontainers-go" */
/* 	"github.com/testcontainers/testcontainers-go/wait" */
/* ) */
/**/
/* const DOCKER_HOST = "host.docker.internal" */
/**/
/* func TestWithRedis(t *testing.T) { */
/* 	ctx := context.Background() */
/* 	ports, err := nat.NewPort("tcp", strconv.Itoa(defaults.REDIS_PORT)) */
/* 	if err != nil { */
/* 		t.Fail() */
/* 		logger.Logger.Error("Unexpected error err=", err) */
/* 	} */
/**/
/* 	req := testcontainers.ContainerRequest{ */
/* 		Image:        "redis:latest", */
/* 		ExposedPorts: []string{string(ports)}, */
/* 		WaitingFor:   wait.ForLog("Ready to accept connections"), */
/* 	} */
/* 	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ */
/* 		ContainerRequest: req, */
/* 		Started:          true, */
/* 	}) */
/* 	if err != nil { */
/* 		logger.Logger.Fatalf("Could not start redis: %s", err) */
/* 	} */
/* 	defer func() { */
/* 		if err := redisC.Terminate(ctx); err != nil { */
/* 			logger.Logger.Fatalf("Could not stop redis: %s", err) */
/* 		} */
/* 	}() */
/**/
/* 	// 1. Create new redis client */
/* 	mappedPort, err := redisC.MappedPort(ctx, ports) */
/* 	addr := fmt.Sprintf("%s:%v", DOCKER_HOST, mappedPort.Int()) */
/* 	client := redis.NewRedisClient(addr) */
/* 	defer client.Close() */
/**/
/* 	// 2. Create new redis adapters handler */
/* 	readAdapter, err := redis.RedisClientAdapterFactory(&model.RepositoryClientConfig{ */
/* 		Name:    "redis", */
/* 		Variant: "redis", */
/* 		Action:  "read", */
/* 	}, client) */
/**/
/* 	writeAdapter, err := redis.RedisClientAdapterFactory(&model.RepositoryClientConfig{ */
/* 		Name:    "redis", */
/* 		Variant: "redis", */
/* 		Action:  "write", */
/* 	}, client) */
/**/
/* 	// 3. Test Read request */
/* 	readCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second) */
/* 	defer cancel() */
/* 	readRes := readAdapter.Call(readCtx) */
/* 	if readRes.Error != goRedis.Nil { // go-redis returns "nil" error when no match */
/* 		t.Fail() */
/* 		logger.Logger.Error(readRes) */
/* 		return */
/* 	} */
/**/
/* 	// 4. Test Write request */
/* 	writeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second) */
/* 	defer cancel() */
/* 	writeRes := writeAdapter.Call(writeCtx) */
/* 	if writeRes.Error != nil { */
/* 		t.Fail() */
/* 		logger.Logger.Error(writeRes.Error) */
/* 		return */
/* 	} */
/* } */
