package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"time"

	"gofiber-starterkit/app/models"
	"gofiber-starterkit/pkg/utils"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
	"golang.org/x/crypto/bcrypt"
)

var Migrations = migrate.NewMigrations()

//go:embed *.sql
var sqlMigrations embed.FS

func init() {
	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}

	if os.Getenv("APP_ENV") == "" || os.Getenv("APP_ENV") == "local" || os.Getenv("APP_ENV") == "development" {
		if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Error().Err(err).Msg("Error loading .env file")
		}
	}
}

func main() {
	ctx := context.Background()

	seed := false
	for _, arg := range os.Args[1:] {
		if arg == "--seed" {
			seed = true
		}
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(utils.DatabaseConnectionString()),
	))

	pg := bun.NewDB(sqldb, pgdialect.New())
	defer pg.Close()

	migrator := migrate.NewMigrator(pg, Migrations)

	if err := migrator.Init(ctx); err != nil {
		fmt.Printf("Failed to initialize migrator: %v\n", err)
		os.Exit(1)
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	if group.IsZero() {
		fmt.Println("No new migrations to run")
	} else {
		fmt.Printf("Migrated to %s\n", group)
	}

	if seed {
		if err := seeder(pg, ctx); err != nil {
			fmt.Printf("Failed to run seeder: %v\n", err)
			os.Exit(1)
		}
	}
}

func seeder(pg *bun.DB, ctx context.Context) error {
	fmt.Println("Running seeder...")

	users := []struct {
		username string
		password string
		email    string
	}{
		{
			username: "admin",
			password: "admin123",
			email:    "admin@example.com",
		},
	}

	for _, u := range users {
		var existing models.User
		err := pg.NewSelect().Model(&existing).Where("username = ?", u.username).Scan(ctx)
		if err == nil {
			fmt.Printf("User %s already exists, skipping\n", u.username)
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", u.username, err)
		}

		hashedPasswordStr := string(hashedPassword)
		user := models.User{
			Username:     u.username,
			PasswordHash: &hashedPasswordStr,
			Email:        u.email,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		_, err = pg.NewInsert().Model(&user).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", u.username, err)
		}

		fmt.Printf("Seeded user: %s\n", u.username)
	}

	fmt.Println("Seeding completed")
	return nil
}
