package database

import (
	"database/sql"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PackageModel struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"index"`
	Latest    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PackageVersionModel struct {
	Id           uuid.UUID `gorm:"primaryKey"`
	Version      string
	Author       string
	Description  sql.NullString
	Dependencies datatypes.JSON

	DistIntegrity    string
	DistShasum       string
	DistFileCount    uint
	DistUnpackedSize uint

	CreatedAt time.Time
}

type AuditLogModel struct {
	Id      uint `gorm:"primaryKey"`
	UserKey string
}

func Open(filePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("sqlite.db?_pragma=journal_mode(WAL)"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&PackageModel{}, &PackageVersionModel{}, &AuditLogModel{}); err != nil {
		return db, err
	}

	return db, nil
}
