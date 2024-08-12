package main

import (
	"context"
	"health/config"
	pb "health/genproto/health_analytics"
	mongoDb "health/mongodb"
	"health/service"
	"log"
	"net"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", config.Load().HEALTH_SERVICE)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	mongoClient, mongodb, err := mongoDb.NewMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis serverining manzili
		Password: "",               // Parol agar mavjud bo'lsa
		DB:       0,                // Default DB ni ishlatish
	})

	mongoDbRepo := mongoDb.NewHealth(mongodb,rdb)
	HelathService := service.NewHealthService(mongoDbRepo)

	server := grpc.NewServer()
	pb.RegisterHealthAnalyticsServiceServer(server, HelathService)

	log.Printf("Server is listening on port %s\n", config.Load().HEALTH_SERVICE)
	if err = server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
