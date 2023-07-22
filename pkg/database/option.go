package database

import (
	"github.com/rasteiro11/PogCore/pkg/logger"
	"gorm.io/gorm"
)

func WithMaxIdleConns(maxIddleCons int) EngineOpt {
	return func(d *gorm.DB) {
		db, err := d.DB()
		if err != nil {
			logger.Global().Errorf("[database.WithMaxIdleConns] d.DB() returned error: %+v\n", err)
			return
		}

		db.SetMaxIdleConns(maxIddleCons)
	}
}

func WithMaxOpenConns(maxOpenConns int) EngineOpt {
	return func(d *gorm.DB) {
		db, err := d.DB()
		if err != nil {
			logger.Global().Errorf("[database.WithMaxOpenConns] d.DB() returned error: %+v\n", err)
			return
		}

		db.SetMaxOpenConns(maxOpenConns)
	}
}
