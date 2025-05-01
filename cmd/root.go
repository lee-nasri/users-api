/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"users-api/pkg/logx"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "Start users application",
	Long:  `Start users application`,
	Run: func(_ *cobra.Command, _ []string) {

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tasks-api.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func gracefullShutdown(quit chan os.Signal, logger *zap.Logger, server *echo.Echo, redis *redis.Client) {
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if server != nil {
		logx.Log().Info("Server is shutting down...")
		if err := server.Shutdown(ctx); err != nil {
			logx.Infof(ctx, "[!] Fail to shutdown service caz: %s", err)
		}
	}

	if redis != nil {
		logx.Log().Info("Redis is shutting down...")
		if err := redis.Close(); err != nil {
			logx.Infof(ctx, "[!] Fail to shutdown redis caz: %s", err)
		}
	}

	_ = logger.Sync()
	os.Exit(0)
}
