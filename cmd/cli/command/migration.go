package command

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"karlota.aasumitro.id/config"
	"karlota.aasumitro.id/internal/utils/security"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var cfg *config.Config

var CliCommand = &cobra.Command{
	Long: `SkilledIn Reference Service CLI App`,
}

var migrateCmd = &cobra.Command{
	Use:  "migration",
	Long: `migrate cmd is used for database migration: migrate < up | down >`,
}

var migrateUpCmd = &cobra.Command{
	Use:  "up",
	Long: `Command to upgrade database migration`,
	Run: func(cmd *cobra.Command, args []string) {
		migration, err := initGoMigrate()
		if err != nil {
			log.Printf("migrate down error: %v \n", err)
			return
		}
		if err := migration.Up(); err != nil {
			log.Printf("migrate up error: %v \n", err)
			return
		}
		log.Println("Migrate up done with success")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:  "down",
	Long: `Command to downgrade database`,
	Run: func(cmd *cobra.Command, args []string) {
		migration, err := initGoMigrate()
		if err != nil {
			log.Printf("migrate down error: %v \n", err)
			return
		}
		if err := migration.Down(); err != nil {
			log.Printf("migrate down error: %v \n", err)
			return
		}
		log.Println("Migrate down done with success")
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:  "version",
	Long: `Command to see database migration version`,
	Run: func(cmd *cobra.Command, args []string) {
		migration, err := initGoMigrate()
		if err != nil {
			log.Printf("migrate down error: %v \n", err)
			return
		}
		version, dirty, err := migration.Version()
		if err != nil {
			log.Printf("migrate up error: %v \n", err)
			return
		}
		log.Printf("Database Version %d is dirty: %s \n",
			version, strconv.FormatBool(dirty))
	},
}

var migrateSeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with initial data",
	Run: func(cmd *cobra.Command, args []string) {
		db := cfg.Infra.SQLPool
		defer func() { _ = db.Close() }()
		users := []string{"ani", "beni", "ceni", "dani", "eri",
			"fani", "gani", "heni", "indra", "jeni", "kiki", "leni",
			"meli", "nani", "oni", "peni", "qori", "reni", "santi",
			"tini", "umar", "vina", "weni", "xena", "yogi", "zeni",
		}
		stmt, err := db.Prepare("INSERT INTO users (display_name, email, password, email_verified_at, " +
			"password_updated_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatalf("failed to prepare statement: %v", err)
		}
		defer func() { _ = stmt.Close() }()

		now := time.Now()
		for _, user := range users {
			h := security.PasswordHash{Stored: "", Supplied: "Test@1234"}
			pwd, err := h.MakePassword(security.Parallelization)
			if err != nil {
				log.Fatalf("failed to insert user %s: %v", user, err)
			}
			email := fmt.Sprintf("%s@karlota.id", user)
			_, err = stmt.Exec(user, email, pwd, now, now, now, now)
			if err != nil {
				log.Printf("failed to insert user %s: %v", user, err)
			}
		}
		log.Println("Seeding completed successfully")
	},
}

func initGoMigrate() (instance *migrate.Migrate, err error) {
	// create file source
	fileSource, err := (&file.File{}).Open("file://db/migrations")
	if err != nil {
		return nil, err
	}
	// create sqlite3 driver
	driver, err := sqlite3.WithInstance(
		cfg.Infra.SQLPool, &sqlite3.Config{})
	if err != nil {
		return nil, err
	}
	// create migrate instance
	instance, err = migrate.NewWithInstance(
		"file", fileSource, "reference", driver)
	if err != nil {
		return nil, err
	}
	// return migrate instance
	return instance, nil
}

func init() {
	viper.SetConfigFile(".env")
	cfg = config.Load().With(
		context.Background(),
		config.SQLiteConnection())
	CliCommand.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateVersionCmd)
	migrateCmd.AddCommand(migrateSeedCmd)
}
