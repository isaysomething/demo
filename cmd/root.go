package cmd

import (
	"log"
	"os"

	"github.com/clevergo/demo/internal/core"
	"github.com/knadh/koanf"
	"github.com/spf13/cobra"
)

var (
	cfg     = core.Config{}
	cfgFile *string
	k       = koanf.New(".")
)

func init() {
	cfgFile = rootCmd.PersistentFlags().StringP("config", "c", "config.toml", "config file")
	parseConfig()
	rootCmd.PersistentFlags().Parse(os.Args[1:])
	rootCmd.AddCommand(
		serveCmd,
		apiCmd,
		migrateCmd,
	)
}

var rootCmd = &cobra.Command{
	Use:   "Demo",
	Short: "Demo is a CleverGo application project template",
	Long:  `Demo is a CleverGo application project template`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute executes commands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed execute: %s", err)
	}
}
