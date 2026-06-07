package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "energy-prometheus-exporter",
		Short: "The homewizard prometheus metrics exporter",
	}

	rootCmd.AddCommand(newHomewizardCommand())
	rootCmd.AddCommand(newServeCommand())
	return rootCmd
}

func initConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetEnvPrefix("")
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
