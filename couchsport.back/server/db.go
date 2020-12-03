package server

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func mustOpenDb() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	host := viper.GetString("DATABASE_HOST")
	user := viper.GetString("DATABASE_USER")
	password := viper.GetString("DATABASE_PASSWORD")
	database := viper.GetString("DATABASE_DATABASE")
	params := viper.GetString("DATABASE_PARAMS")
	port := viper.GetString("DATABASE_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, password, host, port, database, params)
	log.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	_ = os.Setenv("DATABASE_PASSWORD", "")
	_ = os.Setenv("DATABASE_ROOT_PASSWORD", "")

	// db.LogMode(c.Verbose)

	// validations.RegisterCallbacks(db)

	return db
}
