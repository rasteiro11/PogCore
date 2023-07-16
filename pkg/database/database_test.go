package database_test

import (
	"flashcards/core/database"
	"log"
	"testing"

	"gorm.io/gorm"
)

func createDatabase(t *testing.T) database.Database {
	database, err := database.NewDatabase(database.GetMysqlEngineBuilder)
	if err != nil {
		log.Fatalf("[createDatabase] database.NewDatabase() returned error: %+v\n", err)
	}

	return database
}

func getEntities() []any {
	return []any{
		User{},
		Card{},
	}
}

type (
	User struct {
		gorm.Model
		Name string
		Age  int
	}

	Card struct {
		gorm.Model
		Title string
	}
)

func TestNewDatabase(t *testing.T) {
	createDatabase(t)
}

func TestTestMigration(t *testing.T) {
	database := createDatabase(t)

	err := database.Migrate(getEntities()...)
	if err != nil {
		log.Fatalf("[TestTestMigration] database.Migrate() returned error: %+v\n", err)
	}
}

func dropTables(t *testing.T, db database.Database) {
	migrator := db.Conn().Migrator()
	tableNames, err := migrator.GetTables()
	if err != nil {
		log.Fatalf("[dropTables] db.Conn.GetTables() returned error: %+v\n", err)
	}

	for _, tableName := range tableNames {
		err := migrator.DropTable(tableName)
		if err != nil {
			log.Fatalf("[dropTables] db.DropTable() returned error: %+v\n", err)
		}
	}
}

func TestTeardown(t *testing.T) {
	database := createDatabase(t)
	dropTables(t, database)
}
