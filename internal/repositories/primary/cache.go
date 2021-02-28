package primary

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/maxim-kuderko/service-template/pkg/requests"
	"github.com/maxim-kuderko/service-template/pkg/responses"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Cache struct {
	origin Repo
	redis  *redis.ClusterClient
	ttl    time.Duration
}

func NewCache(origin Repo, v *viper.Viper) Repo {
	return &Cache{origin: origin, redis: redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          strings.Split(v.GetString(`PRIMARY_CACHE_REDIS_CLUSTER`), `,`),
		RouteByLatency: true,
	}),
		ttl: time.Millisecond * v.GetDuration(`CACHE_EXPIRATION_MS`),
	}
}
func NewCachedDB(v *viper.Viper) Repo {
	return NewCache(NewDb(v), v)
}

func (c *Cache) Get(r requests.Get) (responses.Get, error) {
	return c.loadOrStore(r)
}

func (c *Cache) loadOrStore(r requests.Get) (responses.Get, error) {
	resp, err := c.load(r)
	if err == nil {
		return resp, err
	}
	resp, err = c.origin.Get(r)
	if err != nil {
		return resp, err
	}
	return resp, c.store(r, resp)
}

const CACHE_KEY = `v1:key:%s`

func (c *Cache) load(r requests.Get) (responses.Get, error) {
	output := responses.Get{}
	resp, err := c.redis.Get(r.Context(), fmt.Sprintf(CACHE_KEY, r.Key)).Bytes()
	if err != nil {
		return output, err
	}

	return output, jsoniter.ConfigFastest.Unmarshal(resp, &output)
}

func (c *Cache) store(r requests.Get, resp responses.Get) error {
	b, err := jsoniter.ConfigFastest.Marshal(resp)
	if err != nil {
		return err
	}
	return c.redis.Set(r.Context(), fmt.Sprintf(CACHE_KEY, r.Key), b, c.ttl).Err()
}
