package cmd

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

func newServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve Prometheus exported metrics",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cmd, "")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			reg := prometheus.NewRegistry()

			c := buildCollectorForHomewizard(cmd)
			if c != nil {
				reg.Register(c)
			}

			http.HandleFunc("/health-check", healthCheck)
			http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
			http.ListenAndServe(":8080", nil)
			return nil
		},
	}
	registerHomewizard(cmd)
	return cmd
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
