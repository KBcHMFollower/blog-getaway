package main

import (
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"test-plate/config"
	"test-plate/internal/app"
	"test-plate/internal/logger"
)

func init() {
	var (
		argConfigPath string
	)

	runCmd := &cobra.Command{
		Use:   "run-gw",
		Short: "Start gateway server",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.MustLoad(argConfigPath)
			log := logger.CongigurateLogger(cfg.Env)
			log.Debug("debug messages are enabled")

			app := app.New(log, cfg)
			app.AddMiddlewares().
				AddHandlers().
				Run()

			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			<-done
			log.Info("stopping server")

			app.Stop()

			log.Info("server stopped")
		},
	}
	runCmd.Flags().StringVar(&argConfigPath, "cfg", "config/local.yaml", "path to config")
	rootCmd.AddCommand(runCmd)
}
