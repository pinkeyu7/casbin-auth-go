package repository

import (
	"casbin-auth-go/internal/system/system"

	"github.com/go-redis/redis/v7"
)

type Cache struct {
	redis *redis.ClusterClient
}

func NewCache(r *redis.ClusterClient) system.Cache {
	return &Cache{redis: r}
}
