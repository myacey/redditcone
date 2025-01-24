package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/myacey/redditclone/internal/apiserver"
	"github.com/myacey/redditclone/internal/logging"
	"github.com/myacey/redditclone/internal/repository/mongorepo"
	"github.com/myacey/redditclone/internal/repository/postgresrepo"
	"github.com/myacey/redditclone/internal/repository/redisrepo"
	"github.com/myacey/redditclone/internal/service"
	"github.com/myacey/redditclone/internal/token/jwttoken"
)

func main() {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "error loading environment variables: %v", err)
		os.Exit(1)
	}

	logger := logging.ConfigureLogger()
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatalf("cant sync logger: %v", err)
		}
	}()

	mongoClient, err := mongorepo.ConfigureMongoClient()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("mongo initialized")

	rdb, err := redisrepo.ConfigureRedisClient()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("redis initialized")
	defer func() {
		err = rdb.Close()
		if err != nil {
			logger.Fatalf("cant close connection to redis: %v", err)
		}
	}()

	postgresDB, err := postgresrepo.ConfigurePostgres()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("postgres initialized")

	tokenMaker := jwttoken.NewJWTToken([]byte(os.Getenv("JWT_SECRET_KEY")))

	service := service.NewService(postgresDB, mongoClient, "redditclone", rdb, tokenMaker, logger)

	server := apiserver.NewServer(logger, service, tokenMaker)
	server.Start()
}
