package cmd

import (
	stdlog "log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     = &Config{}
	k       = koanf.New(".")
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.toml", "config file")

	rootCmd.AddCommand(
		serveCmd,
		migrateCmd,
	)
}

var rootCmd = &cobra.Command{
	Use:   "saysth",
	Short: "SaySth is an fast and flexible CMS",
	Long:  `SaySth is an fast and flexible CMS`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute executes commands.
func Execute() {
	if err := parseConfig(); err != nil {
		stdlog.Fatalf("failed parse config: %s", err)
	}

	if err := rootCmd.Execute(); err != nil {
		stdlog.Fatalf("failed execute: %s", err)
	}
}

func parseConfig() error {
	parser := toml.Parser()
	for _, f := range []string{"configs/config.toml", cfgFile} {
		if err := k.Load(file.Provider(f), parser); err != nil {
			return err
		}
	}

	if err := k.Unmarshal("", cfg); err != nil {
		return err
	}

	return nil
}
