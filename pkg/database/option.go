package database

import (
	"log"

	"gorm.io/gorm"
)

func WithMaxIdleConns(maxIddleCons int) EngineOpt {
	return func(d *gorm.DB) {
		db, err := d.DB()
		if err != nil {
			log.Printf("[database.WithMaxIdleConns] d.DB() returned error: %+v\n", err)
			return
		}

		db.SetMaxIdleConns(maxIddleCons)
	}
}

func WithMaxOpenConns(maxOpenConns int) EngineOpt {
	return func(d *gorm.DB) {
		db, err := d.DB()
		if err != nil {
			log.Printf("[database.WithMaxOpenConns] d.DB() returned error: %+v\n", err)
			return
		}

		db.SetMaxOpenConns(maxOpenConns)
	}
}
