package db

import (
	"fmt"

	"github.com/ezep02/microservicios/internal/auth/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase() (*Database, error) {

	dsn := "host=localhost user=postgres password=7nc4381c4t dbname=auth_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar a la base de datos")
	}

	connection.AutoMigrate(&types.User{})
	connection.AutoMigrate(&types.Session{})
	return &Database{db: connection}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("no se pudo obtener la conexión de base de datos: %w", err)
	}
	return sqlDB.Close()
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}
