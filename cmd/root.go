package cmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "homewizard-prometheus-exporter",
		Short: "The homewizard prometheus metrics exporter",
	}

	rootCmd.AddCommand(newHomewizardCommand())
	rootCmd.AddCommand(newServeCommand())
	return rootCmd
}

func initConfig(cmd *cobra.Command, prefix string) error {
	log.Debugf("%s with prefix: %s", cmd.Name(), prefix)
	v := viper.New()
	v.SetEnvPrefix(prefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	var bindErr error
	visitFlags := func(f *pflag.Flag) {
		if err := v.BindPFlag(f.Name, f); err != nil {
			bindErr = fmt.Errorf("binding flag %q: %w", f.Name, err)
			return
		}
		if !f.Changed && v.IsSet(f.Name) {
			f.Value.Set(fmt.Sprintf("%v", v.Get(f.Name)))
		}
	}

	cmd.Flags().VisitAll(visitFlags)
	if cmd.HasParent() {
		cmd.Parent().PersistentFlags().VisitAll(visitFlags)
	}
	return bindErr
}
