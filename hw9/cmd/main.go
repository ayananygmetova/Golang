package main

import (
	"context"
	"hw9/internal/cache/redis_cache"
	"hw9/internal/http"
	"hw9/internal/store/postgres"
	"log"
)

func main() {
	urlExample := "postgres://demo:123456@localhost:5432/e_store"
	store := postgres.NewDB()
	if err := store.Connect(urlExample); err != nil {
		panic(err)
	}
	defer store.Close()

	cache := redis_cache.NewRedisCache(":8081", 1, 1800)
	//defer cache.Close()

	srv := http.NewServer(context.Background(),
		http.WithAddress(":8081"),
		http.WithStore(store),
		http.WithCache(cache),
	)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()

}
