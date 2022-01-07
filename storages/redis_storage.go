package storages

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/gogotchuri/gocialite"
)

var _ gocialite.GocialStorage = &RedisStorage{}

//Type RedisStorage redis storage for gocialite
type RedisStorage struct {
	client     *redis.Client
	expiration time.Duration
}

//NewRedisStorage returns a new RedisStorage
func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: client,
	}
}

//Get a Gocialite struct from redis
func (s *RedisStorage) Get(key string) (*gocialite.Gocial, error) {
	val, err := s.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	gocial, err := gocialite.Unmarshal([]byte(val))
	if err != nil {
		return nil, err
	}
	return gocial, err
}

//Set a Gocialite struct to redis
func (s *RedisStorage) Set(key string, value *gocialite.Gocial) error {
	val, err := gocialite.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(key, string(val), s.expiration).Err()
}

//Delete a Gocialite struct from redis
func (s *RedisStorage) Delete(key string) error {
	return s.client.Del(key).Err()
}
