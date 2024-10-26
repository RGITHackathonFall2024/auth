package initdb

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/RGITHackathonFall2024/auth/internal/consts"
	"github.com/RGITHackathonFall2024/auth/internal/logic/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDSN() string {
	return fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		os.Getenv(consts.EnvPostgresHost),
		os.Getenv(consts.EnvPostgresUser),
		os.Getenv(consts.EnvPostgresPassword),
		os.Getenv(consts.EnvPostgresDatabase),
		os.Getenv(consts.EnvPostgresPort),
	)
}

func Connect(log *slog.Logger) (*gorm.DB, error) {
	log = log.WithGroup("connect-db")

	dsn := GetDSN()
	log.Debug("Got DSN", slog.String("dsn", dsn))

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Error("Error connecting to database")
		return nil, err
	}

	log.Info("Successfully connected to database")
	return db, nil
}

func InitDB(db *gorm.DB, log *slog.Logger) error {
	log = log.WithGroup("init-db")

	log.Info("Running migrations")
	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Error("Error running migration for User")
		return err
	}

	log.Info("Successfully initialized database")
	return nil
}
