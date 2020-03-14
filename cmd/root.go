package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/spf13/cobra"
)

var (
	cfgFile *string
	cfg     = &Config{}
	k       = koanf.New(".")
)

func init() {
	cfgFile = rootCmd.PersistentFlags().StringP("config", "c", "config.toml", "config file")
	rootCmd.PersistentFlags().Parse(os.Args[1:])
	rootCmd.AddCommand(
		serveCmd,
		migrateCmd,
	)
}

var rootCmd = &cobra.Command{
	Use:   "Demo",
	Short: "Demo is a CleverGo application project template",
	Long:  `Demo is a CleverGo application project template`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return parseConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute executes commands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed execute: %s", err)
	}

}

func parseConfig() error {
	parser := toml.Parser()
	configFS := packr.New("configs", "./../configs")
	// load default configurations.
	configs := configFS.List()
	log.Printf("loading default configurations: %s\n", strings.Join(configs, ", "))
	for _, name := range configs {
		f, err := configFS.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()
		content, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		if err := k.Load(rawbytes.Provider(content), parser); err != nil {
			return err
		}
	}

	if err := k.Load(file.Provider(*cfgFile), parser); err != nil {
		return err
	}

	if err := k.Unmarshal("", cfg); err != nil {
		return err
	}

	return nil
}
