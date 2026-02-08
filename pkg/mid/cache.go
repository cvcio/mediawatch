package mid

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/cvcio/mediawatch/pkg/redis"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type RedisCache struct {
	redis *redis.RedisClient
	log   *logrus.Logger
}

func NewRedisCache(redis *redis.RedisClient, log *logrus.Logger) *RedisCache {
	return &RedisCache{
		redis: redis,
		log:   log,
	}
}

// redis cache middleware
func (c *RedisCache) Cache(ttl time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if ttl == 0 {
				ttl = 24 * time.Hour
			}

			p, err := c.Path(r)
			if err != nil {
				c.log.Error("Error preparing route key")
				return
			}

			k, err := c.Key(r.Body, p)
			if err != nil {
				c.log.Error("Error preparing cache key")
				return
			}
			c.log.Info("Key: ", k)

			v, err := c.Json(k)
			if err != nil {
				return
			}

			if v == nil {
				next.ServeHTTP(w, r)
				return
			}

			c.log.Info("Cache hit", v)

			if len(v) > 0 {
				render.JSON(w, r, v)
				return
			}
		}

		return http.HandlerFunc(fn)
	}
}

func (c *RedisCache) Path(r *http.Request) (string, error) {
	return r.Method + ":" + r.URL.Path + ":" + r.URL.RawQuery, nil
}

func (c *RedisCache) Key(payload interface{}, routeKey string) (string, error) {
	stringPayload := ""
	byteData, err := io.ReadAll(payload.(io.Reader))

	if err != nil {
		return "", err
	}

	stringPayload = string(byteData)
	encodedPayload := base64.StdEncoding.EncodeToString([]byte(stringPayload))
	encodedPayload = routeKey + ":" + encodedPayload

	return encodedPayload, nil
}

func (c *RedisCache) Json(key string) ([]map[string]interface{}, error) {
	res, err := c.redis.Get(key)
	if err != nil {
		c.log.Error("Error getting from cache")
		return nil, nil
	}

	if res == "" {
		return nil, nil
	}

	var j []map[string]interface{}
	if err := json.Unmarshal([]byte(res), &j); err != nil {
		return nil, err
	}

	return j, nil
}

func (c *RedisCache) String(res interface{}) (string, error) {
	s, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(s), nil
}

func (c *RedisCache) Set(r *http.Request, res interface{}) (string, error) {
	p, err := c.Path(r)
	if err != nil {
		c.log.Error("Error preparing route key")
		return "", err
	}

	k, err := c.Key(r.Body, p)
	if err != nil {
		c.log.Error("Error preparing cache key")
		return "", err
	}

	s, err := c.String(res)
	if err != nil {
		c.log.Error("Error stringifying response")
		return "", err
	}

	err = c.redis.Set(k, s, 0)
	if err != nil {
		c.log.Error("Error saving to cache")
		return "", err
	}

	return k, nil
}
