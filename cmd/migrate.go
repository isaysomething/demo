package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

func init() {
	migrateCreateCmd.Flags().String("name", "", "migration name")
	migrateCreateCmd.MarkFlagRequired("name")

	migrateForceCmd.Flags().String("version", "", "version")
	migrateForceCmd.MarkFlagRequired("name")

	migrateCmd.AddCommand(
		migrateDownCmd,
		migrateDropCmd,
		migrateVersionCmd,
		migrateCreateCmd,
		migrateForceCmd,
	)
}

func provideMigrate() *migrate.Migrate {
	db, err := sql.Open(cfg.Migration.Driver, cfg.Migration.DSN)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	source := "file://" + cfg.Migration.Path
	m, err := migrate.NewWithDatabaseInstance(source, cfg.Migration.DB, driver)
	if err != nil {
		log.Fatal(err)
	}

	return m
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "apply all up migrations",
	Long:  `apply all up migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		m := provideMigrate()
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "apply all down migrations.",
	Long:  `apply all down migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		m := provideMigrate()
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
	},
}

var migrateDropCmd = &cobra.Command{
	Use:   "drop",
	Short: "delete everything in database.",
	Long:  `delete everything in database.`,
	Run: func(cmd *cobra.Command, args []string) {
		m := provideMigrate()
		if err := m.Drop(); err != nil {
			log.Fatal(err)
		}
	},
}

var migrateForceCmd = &cobra.Command{
	Use:   "force",
	Short: "sets a migration version.",
	Long:  `sets a migration version.`,
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := strconv.Atoi(cmd.Flag("version").Value.String())
		m := provideMigrate()
		if err := m.Force(version); err != nil {
			log.Fatal(err)
		}
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "display active migration version.",
	Long:  `display active migration version.`,
	Run: func(cmd *cobra.Command, args []string) {
		m := provideMigrate()
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("version: %v, dirty: %t\n", version, dirty)
	},
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a migration",
	Long:  `create a migration`,
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.Flag("name").Value.String()
		now := time.Now().Unix()
		for _, v := range []string{"up", "down"} {
			_, err := os.Create(path.Join(cfg.Migration.Path, fmt.Sprintf("%d_%s.%s.sql", now, name, v)))
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}
