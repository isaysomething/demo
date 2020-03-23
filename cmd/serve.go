package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	serveCmd.AddCommand(
		serveAPICmd,
	)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start front end service",
	Long:  `Start front end service`,
	Run: func(cmd *cobra.Command, args []string) {
		srv, f, err := initializeServer()
		if err != nil {
			log.Fatal(err.Error())
		}
		defer f()

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	},
}
