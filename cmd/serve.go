package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const flagListenAddr = "listen"

// newServeCommand returns the top-level "serve" command. It registers all
// known collectors and exposes them together on a single /metrics endpoint.
func newServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve metrics from all configured devices via HTTP",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			listen, err := cmd.Flags().GetString(flagListenAddr)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			reg := prometheus.NewRegistry()
			err = addSolaredgeCollector(cmd.Flags(), reg, ctx)
			if err != nil {
				return err
			}
			err = addHomewizardCollector(cmd.Flags(), reg)
			if err != nil {
				return err
			}

			return runHTTP(listen, reg)
		},
	}

	cmd.Flags().StringP(flagListenAddr, "l", ":9090", "HTTP listen address")

	// Embed each collector's flags directly on "serve" so the user can configure
	addHomewizardCollectorFlags(cmd.Flags())
	addSolaredgeCollectorFlags(cmd.Flags())

	return cmd
}

// runHTTP starts the HTTP server and blocks until SIGTERM/SIGINT.
func runHTTP(addr string, reg *prometheus.Registry) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health-check", healthCheck)
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	srv := &http.Server{Addr: addr, Handler: mux}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		<-quit
		log.Info("Shutting down…")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	log.Info("http server listening on ", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server: %w", err)
	}
	return nil
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
