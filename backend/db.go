package main

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DB  *pgxpool.Pool
	RDB *redis.Client
	Ctx = context.Background()
)

func initDB() {
	url := mustEnv("DATABASE_URL")
	pool, err := pgxpool.New(Ctx, url)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	DB = pool
	if _, err := DB.Exec(Ctx, mustReadFile("schema.sql")); err != nil {
		log.Fatalf("apply schema: %v", err)
	}
	log.Println("DB ready")
}

func initRedis() {
	addr := mustEnv("REDIS_ADDR")
	RDB = redis.NewClient(&redis.Options{Addr: addr})
	if err := RDB.Ping(Ctx).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}
	log.Println("Redis ready")
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}

func mustReadFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}
