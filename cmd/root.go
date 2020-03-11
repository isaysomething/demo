package cmd

import (
	"io/ioutil"
	"log"

	"github.com/gobuffalo/packr/v2"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
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
		log.Fatalf("failed parse config: %s", err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed execute: %s", err)
	}

}

func parseConfig() error {
	parser := toml.Parser()
	configFS := packr.New("configs", "./../configs")
	// load default configurations.
	for _, name := range configFS.List() {
		log.Printf("loading configuration from %s\n", name)
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

	if err := k.Load(file.Provider(cfgFile), parser); err != nil {
		return err
	}

	if err := k.Unmarshal("", cfg); err != nil {
		return err
	}

	return nil
}
