package cache

import (
	"context"
	"fmt"
	"gin-rest-api-example/internal/config"
	"log"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/suite"
)

type closeFn func() error

type RedisCacheSuite struct {
	suite.Suite
	cacher  Cacher
	closeFn closeFn
}

func TestRedisCache(t *testing.T) {
	suite.Run(t, new(RedisCacheSuite))
}

func (s *RedisCacheSuite) SetupSuite() {
	var err error
	s.cacher, s.closeFn, err = newTestMemoryRedisCacher(s.T())
	s.NoError(err)
}

func (s *RedisCacheSuite) TearDownSuite() {
	if s.closeFn != nil {
		s.closeFn()
	}
}

func (s *RedisCacheSuite) TestFetch() {
	testFetch(s.T(), s.cacher)
}

func newTestMemoryRedisCacher(tb testing.TB) (Cacher, closeFn, error) {
	s := miniredis.RunT(tb)

	conf, _ := config.Load("")
	conf.CacheConfig.Enabled = true
	conf.CacheConfig.Type = "redis"
	conf.CacheConfig.RedisConfig.Cluster = false
	conf.CacheConfig.RedisConfig.Endpoints = []string{s.Addr()}

	cacher, err := newRedisCacher(conf)
	if err != nil {
		tb.Fatalf("failed to create a new redis cacher. err: %v", err)
	}
	return cacher, func() error {
		s.Close()
		return cacher.Close()
	}, nil
}

//func newTestRedisCacher(tb testing.TB) (Cacher, closeFn, error) {
//	tb.Helper()
//	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
//	pool, err := dockertest.NewPool("")
//	if err != nil {
//		tb.Fatalf("failed to connect to docker: %v", err)
//	}
//
//	resource, err := pool.Run("redis", "3.2", nil)
//	if err != nil {
//		tb.Fatalf("could not start resource: %s", err)
//	}
//	closeFn := func() error {
//		if err := pool.Purge(resource); err != nil {
//			tb.Errorf("failed to purge resource: %v", err)
//			return err
//		}
//		return nil
//	}
//	resource.Expire(60 * 5)
//
//	if err = pool.Retry(func() error {
//		db := redis.NewClient(&redis.Options{
//			Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
//		})
//		return db.Ping(context.Background()).Err()
//	}); err != nil {
//		tb.Fatalf("could not connect to docker: %s", err)
//	}
//
//	conf, _ := config.Load("")
//	conf.CacheConfig.Enabled = true
//	conf.CacheConfig.Type = "redis"
//	conf.CacheConfig.RedisConfig.Cluster = false
//	conf.CacheConfig.RedisConfig.Endpoints = []string{fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp"))}
//
//	cacher, err := newRedisCacher(conf)
//	if err != nil {
//		tb.Fatalf("failed to create a new redis cacher. err: %v", err)
//	}
//	return cacher, closeFn, nil
//}

func newTestClusterRedisCacher(tb testing.TB) (Cacher, closeFn, error) {
	tb.Helper()

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		tb.Fatalf("failed to connect to docker: %v", err)
	}
	network, err := pool.Client.CreateNetwork(docker.CreateNetworkOptions{Name: "redis_cluster"})
	if err != nil {
		log.Fatalf("could not create a network to redis cluster: %s", err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "grokzen/redis-cluster",
		Tag:        "6.0.8",
		NetworkID:  network.ID,
		PortBindings: map[docker.Port][]docker.PortBinding{
			"7000/tcp": {{"0.0.0.0", "7000/tcp"}},
			"7001/tcp": {{"0.0.0.0", "7001/tcp"}},
			"7002/tcp": {{"0.0.0.0", "7002/tcp"}},
			"7003/tcp": {{"0.0.0.0", "7003/tcp"}},
			"7004/tcp": {{"0.0.0.0", "7004/tcp"}},
			"7005/tcp": {{"0.0.0.0", "7005/tcp"}},
		},
		Env: []string{
			"IP=0.0.0.0",
			"INITIAL_PORT=7000",
		},
	})
	if err != nil {
		tb.Fatalf("failed to run docker container: %v", err)
	}
	closeFn := func() error {
		if err := pool.Purge(resource); err != nil {
			tb.Errorf("failed to purge resource: %v", err)
			return err
		}
		if err := pool.Client.RemoveNetwork(network.ID); err != nil {
			tb.Errorf("failed to remove %s network: %s", network.Name, err)
			return err
		}
		return nil
	}
	resource.Expire(60 * 5)

	conf, _ := config.Load("")
	conf.CacheConfig.Enabled = true
	conf.CacheConfig.Type = "redis"
	conf.CacheConfig.RedisConfig.Cluster = true
	conf.CacheConfig.RedisConfig.Endpoints = []string{
		"localhost:7000", "localhost:7001", "localhost:7002", "localhost:7003", "localhost:7004", "localhost:7005",
	}
	cli := openRedisCli(conf)
	defer cli.Close()

	if err := pool.Retry(func() error {
		cluster := cli.(*redis.ClusterClient)
		pingErr := cluster.ForEachShard(context.Background(), func(ctx context.Context, client *redis.Client) error {
			result, err := client.Ping(ctx).Result()
			if err != nil {
				return err
			}
			if !strings.EqualFold(result, "pong") {
				return fmt.Errorf("unknown pong message: %s", result)
			}
			return nil
		})
		if pingErr != nil {
			return pingErr
		}

		clusterInfo, infoErr := cluster.ClusterInfo(context.Background()).Result()
		if infoErr != nil {
			return infoErr
		}
		if !strings.Contains(clusterInfo, "cluster_state:ok") {
			return fmt.Errorf("invalid cluster info: %s", clusterInfo)
		}
		return nil
	}); err != nil {
		closeFn()
		tb.Fatalf("failed to connect to redis clusters. err: %v", err)
	}

	cacher, err := newRedisCacher(conf)
	if err != nil {
		tb.Fatalf("failed to create a new redis cacher. err: %v", err)
	}
	return cacher, closeFn, nil
}
