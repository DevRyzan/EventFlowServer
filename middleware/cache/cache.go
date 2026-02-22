package cache

import "time"

// TTL used to set and get cached values data
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
}
