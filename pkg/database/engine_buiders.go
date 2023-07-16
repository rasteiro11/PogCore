package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetMysqlEngineBuilder(dns string, opts ...EngineOpt) (*gorm.DB, error) {
	engine, err := gorm.Open(mysql.Open(dns))
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(engine)
	}

	return engine, nil
}
