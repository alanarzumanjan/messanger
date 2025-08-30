package services

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

func InitDB() {
	url := MustEnv("DATABASE_URL")
	pool, err := pgxpool.New(Ctx, url)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	DB = pool

	schema := mustReadSchema()
	if _, err := DB.Exec(Ctx, schema); err != nil {
		log.Fatalf("apply schema: %v", err)
	}
	log.Println("DB ready")
}

func InitRedis() {
	addr := MustEnv("REDIS_ADDR")
	RDB = redis.NewClient(&redis.Options{Addr: addr})
	if err := RDB.Ping(Ctx).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}
	log.Println("Redis ready")
}

func MustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v
}

func mustReadSchema() string {
	// Сначала пробуем schema.sql, если нет — shema.sql
	if b, err := os.ReadFile("schema.sql"); err == nil {
		return string(b)
	}
	b, err := os.ReadFile("shema.sql")
	if err != nil {
		log.Fatalf("read schema: %v", err)
	}
	return string(b)
}
