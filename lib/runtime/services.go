package runtime

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (R *Runtime) InitServices() error {
	// fmt.Println("InitDB")
	// // open comms to the db
	dia := sqlite.Open(config.Veg.DatabaseConn)

	fmt.Printf("db: %v | %v\n", dia, config.Veg.DatabaseConn)

	db, err := gorm.Open(dia, &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error creating database session service: %w", err)
	}
	R.DB = db

	return nil
}
