package redis_cache

import (
	"context"
	"encoding/json"
	"hw9/internal/cache"
	"hw9/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
)

func (rc RedisCache) Products() cache.ProductsCacheRepo {
	if rc.products == nil {
		rc.products = newProductsRepo(rc.client, rc.expires)
	}

	return rc.products
}

type ProductsRepo struct {
	client  *redis.Client
	expires time.Duration
}

func newProductsRepo(client *redis.Client, exp time.Duration) cache.ProductsCacheRepo {
	return &ProductsRepo{
		client:  client,
		expires: exp,
	}
}

func (c ProductsRepo) Set(ctx context.Context, key string, value []*models.Product) error {
	productsBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err = c.client.Set(ctx, key, productsBytes, c.expires*time.Second).Err(); err != nil {
		return err
	}

	return nil
}

func (c ProductsRepo) Get(ctx context.Context, key string) ([]*models.Product, error) {
	result, err := c.client.Get(ctx, key).Result()
	switch err {
	case nil:
		break
	case redis.Nil:
		return nil, nil
	default:
		return nil, err
	}

	products := make([]*models.Product, 0)
	if err = json.Unmarshal([]byte(result), &products); err != nil {
		return nil, err
	}

	return products, nil
}
