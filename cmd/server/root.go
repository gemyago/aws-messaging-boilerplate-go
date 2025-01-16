package main

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gemyago/aws-messaging-boilerplate-go/internal/app"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/config"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/di"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/diag"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/services"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func newRootCmd(container *dig.Container) *cobra.Command {
	logsOutputFile := ""

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Command to start the server",
	}
	cmd.SilenceUsage = true
	cmd.PersistentFlags().StringP("log-level", "l", "", "Produce logs with given level. Default is env specific.")
	cmd.PersistentFlags().StringVar(
		&logsOutputFile,
		"logs-file",
		"",
		"Produce logs to file instead of stdout. Used for tests only.",
	)
	cmd.PersistentFlags().Bool(
		"json-logs",
		false,
		"Indicates if logs should be in JSON format or text (default)",
	)
	cmd.PersistentFlags().StringP(
		"env",
		"e",
		"",
		"Env that the process is running in.",
	)
	cfg := config.New()
	lo.Must0(cfg.BindPFlags(cmd.PersistentFlags()))
	cmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		rootCtx := cmd.Context()
		err := config.Load(cfg, config.NewLoadOpts().WithEnv(cfg.GetString("env")))
		if err != nil {
			return err
		}

		var logLevel slog.Level
		if err = logLevel.UnmarshalText([]byte(cfg.GetString("defaultLogLevel"))); err != nil {
			return err
		}

		rootLogger := diag.SetupRootLogger(
			diag.NewRootLoggerOpts().
				WithJSONLogs(cfg.GetBool("jsonLogs")).
				WithLogLevel(logLevel).
				WithOptionalOutputFile(logsOutputFile),
		)

		err = errors.Join(
			config.Provide(container, cfg),

			// app layer
			app.Register(container),

			// services
			services.Register(rootCtx, container),

			di.ProvideAll(container,
				di.ProvideValue(rootLogger),
			),
		)
		if err != nil {
			return fmt.Errorf("failed to inject dependencies: %w", err)
		}

		return nil
	}
	return cmd
}
