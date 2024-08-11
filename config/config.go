package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	HEALTH_SERVICE    string
	MongoURI          string
	MongoDBName       string
	RedisAddr         string
	RedisPassword     string
	RedisDB           int
	AUTH_SERVICE_PORT string
}

func Load() Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}

	config := Config{}
	config.HEALTH_SERVICE = cast.ToString(Coalesce("HEALTH_SERVICE", ":50052"))
	config.MongoURI = cast.ToString(Coalesce("MONGO_URI", "mongodb://localhost:27017"))
	config.MongoDBName = cast.ToString(Coalesce("MONGODB_NAME", "health_medicine"))
	config.RedisAddr = cast.ToString(Coalesce("REDIS_ADDR", "localhost:6379"))
	config.RedisPassword = cast.ToString(Coalesce("REDIS_PASSWORD", ""))
	config.RedisDB = cast.ToInt(Coalesce("REDIS_DB", 0))
	config.AUTH_SERVICE_PORT = cast.ToString(Coalesce("AUTH_SERVICE_PORT", ":50051"))

	return config
}

func Coalesce(key string, defaultValue interface{}) interface{} {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
