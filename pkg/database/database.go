package database

import (
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Package struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"index,unique"`
	DistTags  datatypes.JSONMap
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PackageVersion struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	PackageID    uint      `gorm:"index:idversion,unique"`
	Package      Package   `gorm:"foreignKey:PackageID"`
	Version      string    `gorm:"index:idversion,unique"`
	Author       *string
	Description  *string
	Dependencies datatypes.JSONMap
	Engines      datatypes.JSONMap
	Bin          datatypes.JSONMap

	DistIntegrity    string
	DistShasum       string
	DistFileCount    uint
	DistUnpackedSize uint

	CreatedAt time.Time
}

type AuditLog struct {
	ID      uint `gorm:"primaryKey"`
	UserKey string
}

func Open(filePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("sqlite.db?_pragma=journal_mode(WAL)"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&Package{}, &PackageVersion{}, &AuditLog{}); err != nil {
		return db, err
	}

	return db, nil
}
