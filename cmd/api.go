/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"users-api/config"
	"users-api/internal/repository/redis"
	"users-api/internal/rest"
	restUser "users-api/internal/rest/user"
	"users-api/pkg/logx"
	"users-api/pkg/validators"
	"users-api/service/user"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start tasks-api",
	Long:  `Start tasks-api`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize configuration
		config.Init()
		cfg := config.GetConfig()

		if cfg.AppName == "" {
			log.Fatal("Failed to load ENV APP_NAME")
		}

		// Initialize logger
		logger := logx.Init(
			cfg.AppName,
			cfg.AppVersion,
			cfg.AppEnv,
		)

		// Initialize Redis connection
		ctx := context.Background()
		validator, err := validators.NewValidator()
		if err != nil {
			logx.Fatalf(ctx, err, "error initializing validator")
		}

		redisClient, err := redis.NewRedisConnection(ctx, &redis.RedisConfig{
			Host:      cfg.Redis.Host,
			Port:      cfg.Redis.Port,
			Username:  cfg.Redis.Username,
			Password:  cfg.Redis.Password,
			Index:     cfg.Redis.Index,
			KeyPrefix: cfg.Redis.KeyPrefix,
		})
		if err != nil {
			logx.Fatalf(ctx, err, "Failed to initialize Redis connection")
		}

		var (

			// Initialize repository
			repo = redis.NewUserRepository(redisClient, &redis.UserConfig{
				KeyPrefix:    cfg.Redis.KeyPrefix,
				Index:        cfg.Redis.Index,
				DefaultLimit: cfg.Redis.DefaultLimit,
			})

			// Initialize service
			userService = user.NewService(repo)

			// Initialize REST handler
			restHandler = restUser.NewHandler(userService, validator)

			// Initialize and run REST server
			echoServer = rest.NewEchoServer(restHandler)
			e          = echoServer.RunServer(cfg.AppName, cfg.HTTPTimeout)
		)

		go func() {
			err := e.Start(fmt.Sprintf(":%s", cfg.AppPort))
			if err != nil {
				logx.Log().Fatalf("Cannot start server error: %v", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		gracefullShutdown(quit, logger, e, redisClient.Client)
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
