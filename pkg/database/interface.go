package database

import "gorm.io/gorm"

type (
	Database interface {
		Migrate(entities ...any) error
		Conn() *gorm.DB
	}

	EngineBuilder func(dns string, opts ...EngineOpt) (*gorm.DB, error)
	EngineOpt     func(*gorm.DB)
)
