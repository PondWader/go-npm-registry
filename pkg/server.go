package pkg

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/PondWader/go-npm-registry/pkg/config"
	"github.com/PondWader/go-npm-registry/pkg/database"
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

type RequestContext struct {
	DB            *gorm.DB
	Config        config.Config
	Authenticated bool
	UserKey       string
}

func StartServer(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	db, err := database.Open(cfg.DbPath)
	if err != nil {
		return err
	}

	baseReqCtx := RequestContext{
		DB:     db,
		Config: cfg,
	}

	http.HandleFunc("GET /{package}", ContextMiddleware(baseReqCtx, GetPackage))

	return http.ListenAndServe(":8080", nil)
}
