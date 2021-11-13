package main

import (
	"context"

	"hw7/internal/http"
	"hw7/internal/store/postgres"
)

func main() {
	urlExample := "postgres://demo:123456@localhost:5432/e_store"
	store := postgres.NewDB()
	if err := store.Connect(urlExample); err != nil {
		panic(err)
	}
	defer store.Close()

	srv := http.NewServer(context.Background(), ":8080", store)
	if err := srv.Run(); err != nil {
		panic(err)
	}

	srv.WaitForGracefulTermination()

}
