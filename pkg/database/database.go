package database

import (
	"fmt"
	"github.com/rasteiro11/PogCore/pkg/config"
	"gorm.io/gorm"
)

const dnsPattern = "%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local"

type db struct {
	engine *gorm.DB
}

var _ Database = (*db)(nil)

func (d *db) Conn() *gorm.DB {
	return d.engine
}

func (d *db) Migrate(entities ...any) error {
	for _, entity := range entities {
		if err := d.engine.AutoMigrate(entity); err != nil {
			return err
		}
	}

	return nil
}

func NewDatabase(engine EngineBuilder, opts ...EngineOpt) (Database, error) {
	user := config.Instance().RequiredString("DATABASE_USER")
	password := config.Instance().RequiredString("DATABASE_PASSWORD")
	addr := config.Instance().RequiredString("DATABASE_ADDR")
	database := config.Instance().RequiredString("DATABASE")

	dns := fmt.Sprintf(dnsPattern, user, password, addr, database)
	e, err := engine(dns, opts...)
	if err != nil {
		return nil, err
	}

	db := &db{
		engine: e,
	}

	return db, nil
}
